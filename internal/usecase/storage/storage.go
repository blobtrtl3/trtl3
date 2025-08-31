package storage

import (
	"database/sql"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Storage interface {
	Save(bi *domain.BlobInfo, blob *[]byte) (bool, error)
	FindByBucket(bucket string) (*[]domain.BlobInfo, error)
	FindByBucketAndID(bucket string, id string) (*domain.BlobInfo, error)
	Delete(bucket string, id string) (bool, error)
}

type BlobStorage struct {
	db *sql.DB
}

func NewBS(db *sql.DB) Storage {
	return &BlobStorage{
		db: db,
	}
}
