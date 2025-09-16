package jobs

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/repo/signatures"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
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
