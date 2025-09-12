package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blobtrtl3/trtl3/api/routes"
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

	var path = "blobs"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory to save blobs, reason: %s", err)
	}

	storage := storage.NewBlobStorage(conn, path)
	signatures := signatures.NewMapSignatures(signeds)

	routes.NewRoutesCtx(r, storage, signatures).SetupRoutes()

	job := jobs.NewJobs(storage, path, signatures)
	go job.Start(time.Duration(jobInterval) * time.Minute)

	r.Run(":7713")
}
