package db

import (
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb/v2"
)

func NewDbConn() *sql.DB {
	db, err := sql.Open("duckdb", "db")
	if err != nil {
		log.Fatalf("Could not create connection in database, reason: %s", err)
	}

	return db
}
