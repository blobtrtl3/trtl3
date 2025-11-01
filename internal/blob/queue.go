package blob

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/blobtrtl3/trtl3/pkg/domain"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	Info    *domain.BlobInfo
	Blob    io.Reader
	Retries int
}

type Queue struct {
	redis    *redis.Client
	wg       *sync.WaitGroup
	blobRepo *Repository
	ctx      context.Context
}

func NewQueue(br *Repository, r *redis.Client, ctx context.Context) *Queue {
	queue := &Queue{
		redis:    r,
		wg:       &sync.WaitGroup{},
		blobRepo: br,
		ctx:      ctx,
	}

	return queue
}

func (q *Queue) Push(blobInfo *domain.BlobInfo, r io.Reader) error {
	task := Task{
		Info: blobInfo,
		Blob: r,
	}

	if err := q.redis.LPush(q.ctx, "blobs.queue", task).Err(); err != nil {
		return err
	}

	return nil
}

func (q *Queue) SetupWorkers(workers int) {
	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go func(id int) {
			defer q.wg.Done()
			for {
				values, err := q.redis.BLPop(q.ctx, 0, "logs.queue").Result()
				if err != nil {
					continue
				}

				var task Task
				if err := json.Unmarshal([]byte(values[1]), &task); err != nil {
					//retry
					continue
				}

				if _, err := q.blobRepo.Save(task.Info, task.Blob); err != nil {
					fmt.Print(err)
					continue
				}

				fmt.Printf("worker id[%d]", id)
			}
		}(i + 1)
	}
}
