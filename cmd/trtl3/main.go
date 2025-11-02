package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blobtrtl3/trtl3/internal/blob"
	"github.com/blobtrtl3/trtl3/internal/http/router"
	"github.com/blobtrtl3/trtl3/internal/infra"
	"github.com/blobtrtl3/trtl3/internal/jobs"
	"github.com/gin-gonic/gin"
)

// @title Trtl3 API
// @version 1.0
// @description Blob storage api
func main() {
	r := gin.Default()

	conn := infra.NewDbConn()
	defer conn.Close()

	var path = "blobs"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	ctx := context.Background()

	redisClient := infra.NewRedistClient(ctx)

	blobRepo := blob.NewRepository(conn, path)

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

	blobQueue := blob.NewQueue(blobRepo, redisClient, ctx)
	blobQueue.SetupWorkers(workers)

	blobService := blob.NewService(blobRepo, redisClient, blobQueue)

	router.NewRouterCtx(r, blobService, redisClient).SetupRouter()

	job := jobs.NewJobs(blobRepo, path, redisClient)
	go job.Start(time.Duration(jobInterval) * time.Minute) // take interval from env

	r.Run(":7713")
}
