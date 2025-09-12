package handler

import (
	"github.com/blobtrtl3/trtl3/internal/usecase/signatures"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
)

type BlobHandler struct {
	storage storage.Storage
	signatures signatures.Signatures
}

func NewBlob(st storage.Storage, sg signatures.Signatures) *BlobHandler {
	return &BlobHandler{
		storage: st,
		signatures: sg,
	}
}
