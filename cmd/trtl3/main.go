package main

import (
	"log"
	"os"
	"strconv"
	"time"

	handler "github.com/blobtrtl3/trtl3/api/handler/blob"
	"github.com/blobtrtl3/trtl3/api/middleware"
	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/jobs"
	"github.com/blobtrtl3/trtl3/internal/usecase/signatures"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

// @title Trtl3 API
// @version 1.0
// @description Blob storage api
func main() {
	jobInterval, err := strconv.Atoi(os.Getenv("JOB_INTERVAL"))
	if err != nil {
		jobInterval = 5
	}

	r := gin.Default()

	conn := db.NewDbConn()
	defer conn.Close()

	signeds := map[string]domain.Signature{}

	_, err = conn.Exec(`
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

	var path = "blobs"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	storage := storage.NewBlobStorage(conn, path)
	signatures := signatures.NewMapSignatures(signeds)

	blobHandler := handler.NewBlob(storage, signatures)

	protected := r.Group("/blobs", middleware.AuthMiddleware())
	{
		protected.POST("", blobHandler.Save)
		protected.GET("", blobHandler.FindByBucket)
		protected.GET("/:bucket/:id", blobHandler.FindUnique)
		protected.DELETE("/:bucket/:id", blobHandler.Delete)

		protected.GET("/download/:bucket/:id", blobHandler.Download)

		protected.POST("/sign", blobHandler.Sign)
	}

	serve := r.Group("/b", middleware.SignMiddleware(signatures))
	{
		serve.GET("", blobHandler.Serve)
	}

	job := jobs.NewJobs(storage, path, signatures)
	go job.Start(time.Duration(jobInterval) * time.Minute)

	r.Run(":7713")
}
