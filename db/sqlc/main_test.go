package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Yogksai/simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

// *sql.DB is not used in this package, so we use *pgxpool.Pool instead
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: " + err.Error())
	}
	ctx := context.Background()

	testPool, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testPool)
	os.Exit(m.Run())
}
