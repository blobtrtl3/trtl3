package storage

import (
	"database/sql"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Storage interface {
	Save(blobInfo *domain.BlobInfo, blobBytes []byte) (bool, error)
	FindByBucket(bucket string) ([]domain.BlobInfo, error)
	FindByBucketAndID(bucket string, id string) (*domain.BlobInfo, error)
	Delete(bucket string, id string) (bool, error)
	DownloadByID(bucket string, id string) ([]byte, error)
	FindAll() ([]domain.BlobInfo, error)
}

type BlobStorage struct {
	db *sql.DB
}

func NewBlobStorage(db *sql.DB) Storage {
	return &BlobStorage{
		db: db,
	}
}
