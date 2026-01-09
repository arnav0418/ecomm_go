package main

import (
	"log"

	"github.com/arnav0418/ecomm_go/db"
	"github.com/arnav0418/ecomm_go/ecomm-api/server"
	"github.com/arnav0418/ecomm_go/ecomm-api/storer"
)

func main() {
	db, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	log.Printf("Successfully connected to db")

	// doing something with the database
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
}