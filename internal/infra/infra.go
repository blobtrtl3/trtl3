package infra

import (
	"context"
	"database/sql"
	"log"

	"github.com/blobtrtl3/trtl3/pkg/domain"
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

type SignaturesCache interface {
	Set(key string, val domain.Signature)
	Get(key string) *domain.Signature
	Delete(key string)
	FindAll() []string
}

type MemSignaturesCache map[string]domain.Signature

func NewMemSignaturesCache() SignaturesCache {
	return MemSignaturesCache{}
}

func (ms MemSignaturesCache) Set(key string, val domain.Signature) {
	ms[key] = val
}

func (ms MemSignaturesCache) Get(key string) *domain.Signature {
	if sig, ok := ms[key]; ok {
		return &sig
	}
	return nil
}

func (ms MemSignaturesCache) Delete(key string) {
	delete(ms, key)
}

func (ms MemSignaturesCache) FindAll() []string {
	keys := make([]string, 0, len(ms))
	for key := range ms {
		keys = append(keys, key)
	}
	return keys
}
