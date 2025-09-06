package storage

import (
	"fmt"
	"os"
)

func (bs *BlobStorage) DownloadByID(bucket string, id string) ([]byte, error) {
	blob, err := os.ReadFile(fmt.Sprintf("/tmp/blobs/%s_%s", bucket, id))
	if err != nil {
		return nil, err
	}

	return blob, nil
}
