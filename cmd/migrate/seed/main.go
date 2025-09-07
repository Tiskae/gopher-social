package main

import (
	"log"

	"github.com/tiskae/go-social/internal/db"
	"github.com/tiskae/go-social/internal/env"
	"github.com/tiskae/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	err = db.Seed(store)
	log.Fatal(err)
}
