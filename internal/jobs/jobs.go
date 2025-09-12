package jobs

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/usecase/signatures"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type Jobs struct {
	storage storage.Storage
	signatures signatures.Signatures
	dir     string
}

func NewJobs(storage storage.Storage, dir string, signatures signatures.Signatures) *Jobs {
	return &Jobs{storage: storage, dir: dir, signatures: signatures}
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

