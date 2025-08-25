package main

import (
	"log"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	application := application{
		config: cfg,
	}

	mux := application.mount()

	log.Fatal(application.run(mux))
}
