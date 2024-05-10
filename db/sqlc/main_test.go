package sqlc

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // not used but still needed
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:Anaana123@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer conn.Close()

	testQueries = New(conn)
	os.Exit(m.Run())

}
