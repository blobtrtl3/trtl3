package worker

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type Worker struct {
	storage storage.Storage
	dir string
}

func NewWorker(storage storage.Storage, dir string) *Worker {
	return &Worker{storage: storage, dir: dir}
}

func (w *Worker) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		w.cleanOrphans()
	}
}

func (w *Worker) cleanOrphans() {
	fmt.Println("cleanando")
	blobsinfos, err := w.storage.FindAll()
	if err != nil {
		log.Printf("[worker] error while finding blobs infos: %s", err)
		return
	}

	for _, bi := range blobsinfos {
		path := filepath.Join(w.dir, fmt.Sprintf("%s_%s", bi.Bucket, bi.ID))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("[worker] orphan found (id: %s, bucket: %s)", bi.ID, bi.Bucket)

			_, err := w.storage.Delete(bi.Bucket, bi.ID)
			if err != nil {
				log.Printf("[worker] error deleting orphan: %s", err)
			}
		}
	}
}

