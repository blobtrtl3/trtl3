package jobs

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/repo/signatures"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
	"github.com/blobtrtl3/trtl3/shared"
)

type Jobs struct {
	storage    storage.Storage
	signatures signatures.Signatures
	dir        string
}

func NewJobs(st storage.Storage, dir string, s signatures.Signatures) *Jobs {
	return &Jobs{storage: st, dir: dir, signatures: s}
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
	blobsinfos, err := j.storage.FindAll()
	if err != nil {
		log.Printf("[job] error while finding blobs infos: %s", err)
		return
	}

	for _, bi := range blobsinfos {
		path := filepath.Join(j.dir, shared.GenBlobName(bi.Bucket, bi.ID))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("[job] orphan found (id: %s, bucket: %s)", bi.ID, bi.Bucket)

			_, err := j.storage.Delete(bi.Bucket, bi.ID)
			if err != nil {
				log.Printf("[job] error deleting orphan: %s", err)
			}
		}
	}
}

func (j *Jobs) cleanSignatures() {
	for _, key := range j.signatures.FindAll() {
		if j.signatures.Get(key).TTL.Compare(time.Now()) <= 0 {
			j.signatures.Delete(key)
		}
	}
}
