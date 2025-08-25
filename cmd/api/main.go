package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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

	application := application{
		config: cfg,
	}

	mux := application.mount()

	log.Fatal(application.run(mux))
}
