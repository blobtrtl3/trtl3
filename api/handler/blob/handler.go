package handler

import (
	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type BlobHandler struct {
	storage storage.Storage
	hashmap map[string]domain.Signature
}

func NewBlob(s storage.Storage, hm map[string]domain.Signature) *BlobHandler {
	return &BlobHandler{
		storage: s,
		hashmap: hm,
	}
}
