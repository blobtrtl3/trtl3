package queue

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bq *BlobQueue) Append(blobInfo *domain.BlobInfo, blobBytes []byte) error {
	task := BlobTask{
		Info:  blobInfo,
		Bytes: blobBytes,
	}

	bq.wg.Add(1)
	bq.queue <- task

	return nil
}
