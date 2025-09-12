package jobs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/shared"
)

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
