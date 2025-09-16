package handler

import (
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/internal/repo/signatures"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
)

type BlobHandler struct {
	storage    storage.Storage
	signatures signatures.Signatures
	bloQueue   queue.BlobQueue
}

func NewBlob(
	st storage.Storage,
	sg signatures.Signatures,
	bq queue.BlobQueue,
) *BlobHandler {
	return &BlobHandler{
		storage:    st,
		signatures: sg,
		bloQueue:   bq,
	}
}
