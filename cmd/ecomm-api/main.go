package main

import (
	"log"

	"github.com/arnav0418/ecomm_go/db"
	"github.com/arnav0418/ecomm_go/ecomm-api/handler"
	"github.com/arnav0418/ecomm_go/ecomm-api/server"
	"github.com/arnav0418/ecomm_go/ecomm-api/storer"
	"github.com/ianschenck/envflag"
)

const minSecretKeySize = 32

func main() {
	var secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for jwt signing")
	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be atleast %d characters", minSecretKeySize)
	}

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to database")

	// do something with the database
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}


// main.go ends here