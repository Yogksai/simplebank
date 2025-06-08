package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

// *sql.DB is not used in this package, so we use *pgxpool.Pool instead
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()

	testPool, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testPool)
	os.Exit(m.Run())
}
