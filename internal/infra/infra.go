package infra

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/redis/go-redis/v9"
)

func NewDbConn() *sql.DB {
	db, err := sql.Open("duckdb", "db")
	if err != nil {
		log.Fatalf("Could not create connection in database, reason: %s", err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS blobsinfo (
      id TEXT NOT NULL,
      bucket TEXT NOT NULL,
      mime TEXT NOT NULL,
      size INTEGER NOT NULL,
      created_at TIMESTAMP,
			PRIMARY KEY (id, bucket)
    )
	`)
	if err != nil {
		log.Fatalf("Could not create database table, reason: %s", err)
	}

	return db
}

func NewRedistClient(ctx context.Context) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

