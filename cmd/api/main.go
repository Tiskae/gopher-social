package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tiskae/go-social/internal/store"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT");

	cfg := config{
		addr: PORT,
	}

	storage := store.NewStorage(nil)

	application := application{
		config: cfg,
		store: storage,
	}

	mux := application.mount()

	log.Fatal(application.run(mux))
}
