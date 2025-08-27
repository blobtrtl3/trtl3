package storage

import (
	"database/sql"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Storage interface {
	Save(bi *domain.BlobInfo, blob *[]byte) error
	FindByID(id string) (*domain.BlobInfo, error)
	FindByBucket(bucket string) ([]*domain.BlobInfo, error)
	Delete(id string) error
}

type BlobStorage struct {
	db *sql.DB
}

func NewBS(db *sql.DB) Storage {
	return &BlobStorage{
		db: db,
	}
}

func (bs *BlobStorage) Save(bi *domain.BlobInfo, blob *[]byte) error {
	return nil
}

func (bs *BlobStorage) FindByID(id string) (*domain.BlobInfo, error) {
	bi := domain.BlobInfo{}

	return &bi, nil
}

func (bs *BlobStorage) FindByBucket(bucket string) ([]*domain.BlobInfo, error) {
	return []*domain.BlobInfo{}, nil
}

func (bs *BlobStorage) Delete(id string) error {
	return nil
}
