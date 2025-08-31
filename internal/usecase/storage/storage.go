package storage

import (
	"database/sql"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Storage interface {
	Save(bi *domain.BlobInfo, blob *[]byte) (bool, error)
	FindByID(id string) (*domain.BlobInfo, error)
	FindByBucket(bucket string) (*[]domain.BlobInfo, error)
	Delete(id string) (bool, error)
}

type BlobStorage struct {
	db *sql.DB
}

func NewBS(db *sql.DB) Storage {
	return &BlobStorage{
		db: db,
	}
}
