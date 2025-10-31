package jobs

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/blob"
	"github.com/blobtrtl3/trtl3/internal/infra"
	"github.com/blobtrtl3/trtl3/internal/shared"
)

type Jobs struct {
	blobRepo      *blob.Repository
	signaturesCache infra.SignaturesCache
	dir             string
}

func NewJobs(br *blob.Repository, dir string, sc infra.SignaturesCache) *Jobs {
	return &Jobs{blobRepo: br, dir: dir, signaturesCache: sc}
}

func (j *Jobs) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		j.cleanOrphans()
		j.cleanSignatures()
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

func (j *Jobs) cleanSignatures() {
	for _, key := range j.signaturesCache.FindAll() {
		if j.signaturesCache.Get(key).TTL.Compare(time.Now()) <= 0 {
			j.signaturesCache.Delete(key)
		}
	}
}
