package main

import (
	"log"
	"os"

	"github.com/blobtrtl3/trtl3/api/handler"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	conn := db.NewDbConn()
	defer conn.Close()

	_, err := conn.Exec(`
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

	if err := os.MkdirAll("/tmp/blobs", os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	st := storage.NewBS(conn)

	bh := handler.NewBlob(st)

	r.POST("/blobs", bh.Save)
	r.GET("/blobs", bh.FindByBucketOrID)
	r.DELETE("/blobs", bh.Delete)

	r.Run()
}
