package main

import (
	"github.com/joho/godotenv"
	"github.com/tiskae/go-social/internal/db"
	"github.com/tiskae/go-social/internal/env"
	"github.com/tiskae/go-social/internal/store"
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
		env:     env.GetString("ENV", "development"),
		version: VERSION,
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

	storage := store.NewStorage(db)

	application := application{
		config: cfg,
		store:  storage,
		logger: logger,
	}

	mux := application.mount()

	logger.Fatal(application.run(mux))
}
