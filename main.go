package main

import (
	"context"

	"github.com/Yogksai/simplebank/api"
	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(serverAddress); err != nil {
		panic("cannot start server: " + err.Error())
	}

}
