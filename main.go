package main

import (
	"database/sql"
	"github.com/dangquyit/go-simplebank/api"
	db "github.com/dangquyit/go-simplebank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server")
		return
	}
}
