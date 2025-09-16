package queue

import (
	"sync"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
)

type BlobTask struct {
	Info    *domain.BlobInfo
	Bytes   []byte
	Retries int
}

type BlobQueue struct {
	queue   chan BlobTask
	wg      *sync.WaitGroup
	storage storage.Storage
}

func NewBlobQueue(workers int, s storage.Storage) *BlobQueue {
	q := &BlobQueue{
		queue:   make(chan BlobTask, 24),
		wg:      &sync.WaitGroup{},
		storage: s,
	}

	for range workers {
		go q.worker()
	}

	return q
}
