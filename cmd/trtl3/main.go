package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	handler "github.com/blobtrtl3/trtl3/api/handler/blob"
	"github.com/blobtrtl3/trtl3/api/middleware"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/blobtrtl3/trtl3/internal/worker"
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
	
	var path = filepath.Join(os.TempDir() + "/blobs")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	storage := storage.NewBlobStorage(conn)

	blobHandler := handler.NewBlob(storage)

	protected := r.Group("/blobs", middleware.AuthMiddleware())

	protected.POST("", blobHandler.Save)
	protected.GET("", blobHandler.FindByBucket)
	protected.GET("/:bucket/:id", blobHandler.FindByBucketAndID)
	protected.DELETE("/:bucket/:id", blobHandler.Delete)

	protected.GET("/download/:bucket/:id", blobHandler.DownloadByID)

	r.Static("/b", path)

	worker := worker.NewWorker(storage, path)
	go worker.Start(5 * time.Minute)

	r.Run(":7713")
}
