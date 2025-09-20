package storage

import (
	"database/sql"
	"io"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Storage interface {
	Save(blobInfo *domain.BlobInfo, r io.Reader) (bool, error)
	FindByBucket(bucket string) ([]domain.BlobInfo, error)
	FindUnique(bucket string, id string) (*domain.BlobInfo, error)
	Delete(bucket string, id string) (bool, error)
	Download(bucket string, id string) ([]byte, error)
	FindAll() ([]domain.BlobInfo, error)
}

type BlobStorage struct {
	db  *sql.DB
	dir string
}

func NewBlobStorage(db *sql.DB, dir string) Storage {
	return &BlobStorage{
		db:  db,
		dir: dir,
	}
}
