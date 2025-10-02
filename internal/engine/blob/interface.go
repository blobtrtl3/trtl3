package blob

import (
	"database/sql"
	"io"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type BlobEngine interface {
	Save(blobInfo *domain.BlobInfo, r io.Reader) (bool, error)
	FindByBucket(bucket string) ([]domain.BlobInfo, error)
	FindUnique(bucket string, id string) (*domain.BlobInfo, error)
	Delete(bucket string, id string) (bool, error)
	Download(bucket string, id string) ([]byte, error) // TODO: not load in mem
	FindAll() ([]domain.BlobInfo, error)
}

type BlobEngineImpl struct {
	db  *sql.DB
	dir string
}

func NewBlobEngine(db *sql.DB, dir string) BlobEngine {
	return &BlobEngineImpl{
		db:  db,
		dir: dir,
	}
}
