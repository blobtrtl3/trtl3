package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blobtrtl3/trtl3/internal/engine"
	"github.com/blobtrtl3/trtl3/internal/http/router"
	"github.com/blobtrtl3/trtl3/internal/infra/cache"
	"github.com/blobtrtl3/trtl3/internal/infra/db"
	"github.com/blobtrtl3/trtl3/internal/jobs"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/internal/service"
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

	blobEngine := engine.NewBlobEngine(conn, path)
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

	jobIntervalStr := os.Getenv("JOB_INTERVAL")
	jobInterval := 10 // default value
	if jobIntervalStr != "" {
		if j, err := strconv.Atoi(jobIntervalStr); err == nil && j > 0 {
			jobInterval = j
		} else {
			log.Printf("Invalid JOB_INTERVAL value '%s', using default: %d", jobIntervalStr, jobInterval)
		}
	}

	blobQueue := queue.NewBlobQueue(workers, blobEngine)

	blobService := service.NewBlobService(blobEngine, signaturesCache, blobQueue)

	router.NewRouterCtx(r, blobService, signaturesCache).SetupRouter()

	job := jobs.NewJobs(blobEngine, path, signaturesCache)
	go job.Start(time.Duration(jobInterval) * time.Minute) // take interval from env

	r.Run(":7713")
}
