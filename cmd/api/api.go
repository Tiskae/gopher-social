package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tiskae/go-social/docs" // This is required to generate Swagger docs
	"github.com/tiskae/go-social/internal/auth"
	"github.com/tiskae/go-social/internal/env"
	"github.com/tiskae/go-social/internal/mailer"
	"github.com/tiskae/go-social/internal/ratelimiter"
	"github.com/tiskae/go-social/internal/store"
	"github.com/tiskae/go-social/internal/store/cache"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	apiURL      string
	db          dbConfig
	env         string
	version     string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	expiry time.Duration
	issuer string
}

type sendgridConfig struct {
	apiKey string
}

type mailConfig struct {
	exp       time.Duration
	fromEmail string
	sendgrid  sendgridConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(app.RateLimiterMiddleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Route("/v1", func(r chi.Router) {
		r.
			// With(app.BasicAuthenticationMiddleware()).
			Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // excluded for /comments route internally
				r.Get("/", app.getPostByIDHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))

				r.Route("/comments", func(r chi.Router) {
					r.Get("/", app.getCommentsByPostIDHandler)
					r.Post("/", app.createPostCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				// r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserByIDHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)
	return nil
}
