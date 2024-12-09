package main

import (
	"context"
	"log"

	"main/api"
	db "main/db/sqlc"
	"main/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create a new server", err)
	}

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot run a new server", err)
	}
}
