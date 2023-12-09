package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

// Make Queries object that connect to DB. Can used to do queries
func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err.Error())
	}

	err = testDb.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
