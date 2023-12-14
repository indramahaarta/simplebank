package main

import (
	"database/sql"
	"log"

	"github.com/indramhrt/simplebank/api"
	db "github.com/indramhrt/simplebank/db/sqlc"
	"github.com/indramhrt/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't Load Config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't Connect to DB", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't Start Server", err)
	}
}
