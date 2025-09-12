package jobs

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type Jobs struct {
	storage storage.Storage
	dir     string
}

func NewJobs(storage storage.Storage, dir string) *Jobs {
	return &Jobs{storage: storage, dir: dir}
}

func (j *Jobs) Start(interval time.Duration, m map[string]domain.Signature) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		j.cleanOrphans()
		j.cleanSignatures(m)
	}
}

