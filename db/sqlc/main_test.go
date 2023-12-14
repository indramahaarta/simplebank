package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/indramhrt/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

// Make Queries object that connect to DB. Can used to do queries
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Can't Load Config")
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't Connect to DB", err.Error())
	}

	err = testDb.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
