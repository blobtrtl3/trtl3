package blob

import (
	"io"
	"log"
	"sync"

	"github.com/blobtrtl3/trtl3/pkg/domain"
)

type Task struct {
	Info    *domain.BlobInfo
	Blob    io.Reader
	Retries int
}

type Queue struct {
	queue      chan Task
	wg         *sync.WaitGroup
	blobRepo *Repository
}

func NewQueue(workers int, br *Repository) *Queue {
	q := &Queue{
		queue:      make(chan Task, 24),
		wg:         &sync.WaitGroup{},
		blobRepo: br,
	}

	for range workers {
		go q.worker()
	}

	return q
}


func (q *Queue) Append(blobInfo *domain.BlobInfo, r io.Reader) error {
	task := Task{
		Info: blobInfo,
		Blob: r,
	}

	q.wg.Add(1)
	q.queue <- task

	return nil
}

func (q *Queue) worker() {
	for task := range q.queue {
		_, err := q.blobRepo.Save(task.Info, task.Blob)
		if err != nil {
			task.Retries++
			if task.Retries <= 3 {
				log.Printf("retrying to save blob %s (attempt %d)", task.Info.ID, task.Retries)
				q.Append(task.Info, task.Blob)
			} else {
				log.Printf("failed to save blob %s after %d attempts: %s", task.Info.ID, task.Retries, err)
			}
		}

		q.wg.Done()
	}
}
