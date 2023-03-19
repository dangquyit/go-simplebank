package main

import (
	"database/sql"
	"github.com/dangquyit/go-simplebank/api"
	db "github.com/dangquyit/go-simplebank/db/sqlc"
	"github.com/dangquyit/go-simplebank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server")
		return
	}
}
