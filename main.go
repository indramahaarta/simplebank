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
		log.Fatal("can't load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(&config, store)
	if err != nil {
		log.Fatal("can't create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}
