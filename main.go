package main

import (
	"context"
	"log"

	"github.com/Yogksai/simplebank/api"
	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/Yogksai/simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: " + err.Error())
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(config.ServerAddress); err != nil {
		panic("cannot start server: " + err.Error())
	}

}
