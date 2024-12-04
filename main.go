package main

import (
	"context"
	"log"
	"main/api"

	db "main/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://admin:adminSecret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Run(serverAddress)
	if err != nil {
		log.Fatal("cannot run a new server", err)
	}
}
