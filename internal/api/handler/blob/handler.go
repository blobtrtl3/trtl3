package handler

import (
	"github.com/blobtrtl3/trtl3/internal/cache"
	"github.com/blobtrtl3/trtl3/internal/engine"
	"github.com/blobtrtl3/trtl3/internal/queue"
)

type BlobHandler struct {
	blobEngine      engine.BlobEngine
	signaturesCache cache.SignaturesCache
	bloQueue        queue.BlobQueue
}

func NewBlob(
	be engine.BlobEngine,
	sc cache.SignaturesCache,
	bq queue.BlobQueue,
) *BlobHandler {
	return &BlobHandler{
		blobEngine:      be,
		signaturesCache: sc,
		bloQueue:        bq,
	}
}
