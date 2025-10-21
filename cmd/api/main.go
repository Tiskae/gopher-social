package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/tiskae/go-social/internal/auth"
	"github.com/tiskae/go-social/internal/db"
	"github.com/tiskae/go-social/internal/env"
	"github.com/tiskae/go-social/internal/mailer"
	"github.com/tiskae/go-social/internal/store"
	"github.com/tiskae/go-social/internal/store/cache"
	"go.uber.org/zap"
)

const VERSION = "0.0.1"

//	@title			GopherSocial  API
//	@version		1.0
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Load env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	PORT := env.GetString("PORT", ":8080")

	cfg := config{
		addr:   PORT,
		apiURL: env.GetString("API_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env:     env.GetString("ENV", "development"),
		version: VERSION,
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", "info@gophersocial.com"),
			sendgrid: sendgridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "123"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_AUTH_TOKEN", "example"),
				expiry: time.Hour * 24 * 3, // 3 days
				issuer: "gophersocial",
			},
		},
	}

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database has connected")

	cache := cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)

	storage := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	application := application{
		config:        cfg,
		store:         storage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := application.mount()

	logger.Fatal(application.run(mux))
}
