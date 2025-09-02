package main

import (
	"log"
	"os"

	handler "github.com/blobtrtl3/trtl3/api/handler/blob"
	"github.com/blobtrtl3/trtl3/api/middleware"
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

	storage := storage.NewBS(conn)

	blobHandler := handler.NewBlob(storage)

	protected := r.Group("/blobs", middleware.AuthMiddleware())

	protected.POST("", blobHandler.Save)
	protected.GET("", blobHandler.FindByBucket)
	protected.GET("/:bucket/:id", blobHandler.FindByBucketAndID)
	protected.DELETE("/:bucket/:id", blobHandler.Delete)

	protected.GET("/blobs/download/:id", blobHandler.DownloadByID)

	r.Static("/b", "/tmp/blobs")

	r.Run(":7713")
}
