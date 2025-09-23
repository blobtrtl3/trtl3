package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blobtrtl3/trtl3/internal/api/routes"
	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/jobs"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/internal/repo/signatures"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
	"github.com/gin-gonic/gin"
)

// @title Trtl3 API
// @version 1.0
// @description Blob storage api
func main() {
	r := gin.Default()

	conn := db.NewDbConn()
	defer conn.Close()

	signeds := map[string]domain.Signature{}

	var path = "blobs"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	storage := storage.NewBlobStorage(conn, path)
	signatures := signatures.NewMapSignatures(signeds)

	// Get number of workers from environment variable, default to 10
	workersStr := os.Getenv("WORKERS")
	workers := 10 // default value
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil && w > 0 {
			workers = w
		} else {
			log.Printf("Invalid WORKERS value '%s', using default: %d", workersStr, workers)
		}
	}

	blobQueue := queue.NewBlobQueue(workers, storage)

	routes.NewRoutesCtx(r, storage, signatures, *blobQueue).SetupRoutes()

	job := jobs.NewJobs(storage, path, signatures)
	go job.Start(5 * time.Minute)

	r.Run(":7713")
}
