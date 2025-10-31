package queue

import (
	"io"

	"github.com/blobtrtl3/trtl3/pkg/domain"
)

func (bq *BlobQueue) Append(blobInfo *domain.BlobInfo, r io.Reader) error {
	task := BlobTask{
		Info: blobInfo,
		Blob: r,
	}

	bq.wg.Add(1)
	bq.queue <- task

	return nil
}
