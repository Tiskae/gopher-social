package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tiskae/go-social/internal/db"
	"github.com/tiskae/go-social/internal/env"
	"github.com/tiskae/go-social/internal/store"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := env.GetString("PORT", ":8080")

	cfg := config{
		addr: PORT,
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	storage := store.NewStorage(db)

	application := application{
		config: cfg,
		store:  storage,
	}

	mux := application.mount()

	log.Fatal(application.run(mux))
}
