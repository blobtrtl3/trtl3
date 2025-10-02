package jobs

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/engine"
	"github.com/blobtrtl3/trtl3/internal/infra/cache"
	"github.com/blobtrtl3/trtl3/shared"
)

type Jobs struct {
	blobEngine      engine.BlobEngine
	signaturesCache cache.SignaturesCache
	dir             string
}

func NewJobs(be engine.BlobEngine, dir string, sc cache.SignaturesCache) *Jobs {
	return &Jobs{blobEngine: be, dir: dir, signaturesCache: sc}
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
	blobsinfos, err := j.blobEngine.FindAll()
	if err != nil {
		log.Printf("[job] error while finding blobs infos: %s", err)
		return
	}

	for _, bi := range blobsinfos {
		path := filepath.Join(j.dir, shared.GenBlobName(bi.Bucket, bi.ID))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("[job] orphan found (id: %s, bucket: %s)", bi.ID, bi.Bucket)

			_, err := j.blobEngine.Delete(bi.Bucket, bi.ID)
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
