package queue

import (
	"io"
	"sync"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/engine/blob"
)

type BlobTask struct {
	Info    *domain.BlobInfo
	Blob   	io.Reader
	Retries int
}

type BlobQueue struct {
	queue   chan BlobTask
	wg      *sync.WaitGroup
	blobEngine blob.BlobEngine
}

func NewBlobQueue(workers int, be blob.BlobEngine) *BlobQueue {
	q := &BlobQueue{
		queue:   make(chan BlobTask, 24),
		wg:      &sync.WaitGroup{},
		blobEngine: be,
	}

	for range workers {
		go q.worker()
	}

	return q
}
