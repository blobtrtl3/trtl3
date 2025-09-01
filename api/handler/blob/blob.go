package handler

import (
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type BlobHandler struct {
	storage storage.Storage
}

func NewBlob(s storage.Storage) *BlobHandler {
	return &BlobHandler{
		storage: s,
	}
}

