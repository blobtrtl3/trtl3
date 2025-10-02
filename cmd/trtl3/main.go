package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blobtrtl3/trtl3/internal/api/routes"
	"github.com/blobtrtl3/trtl3/internal/cache"
	"github.com/blobtrtl3/trtl3/internal/engine/blob"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/jobs"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/gin-gonic/gin"
)

// @title Trtl3 API
// @version 1.0
// @description Blob storage api
func main() {
	r := gin.Default()

	conn := db.NewDbConn()
	defer conn.Close()

	var path = "blobs"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	blobEngine := blob.NewBlobEngine(conn, path)
	signaturesCache := cache.NewMemSignaturesCache()

	workersStr := os.Getenv("WORKERS")
	workers := 10 // default value
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil && w > 0 {
			workers = w
		} else {
			log.Printf("Invalid WORKERS value '%s', using default: %d", workersStr, workers)
		}
	}

	blobQueue := queue.NewBlobQueue(workers, blobEngine)

	routes.NewRoutesCtx(r, blobEngine, signaturesCache, *blobQueue).SetupRoutes()

	job := jobs.NewJobs(blobEngine, path, signaturesCache)
	go job.Start(5 * time.Minute) // take interval from env

	r.Run(":7713")
}
