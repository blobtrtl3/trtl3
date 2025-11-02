package jobs

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/blob"
	"github.com/blobtrtl3/trtl3/internal/shared"
	"github.com/redis/go-redis/v9"
)

type Jobs struct {
	blobRepo        *blob.Repository
	redis *redis.Client
	dir             string
	ctx context.Context
}

func NewJobs(br *blob.Repository, dir string, r *redis.Client) *Jobs {
	return &Jobs{blobRepo: br, dir: dir, redis: r}
}

func (j *Jobs) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		j.cleanOrphans()
	}
}

func (j *Jobs) cleanOrphans() {
	blobsinfos, err := j.blobRepo.FindAll()
	if err != nil {
		log.Printf("[job] error while finding blobs infos: %s", err)
		return
	}

	for _, bi := range blobsinfos {
		path := filepath.Join(j.dir, shared.GenBlobName(bi.Bucket, bi.ID))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("[job] orphan found (id: %s, bucket: %s)", bi.ID, bi.Bucket)

			_, err := j.blobRepo.Delete(bi.Bucket, bi.ID)
			if err != nil {
				log.Printf("[job] error deleting orphan: %s", err)
			}
		}
	}
}

