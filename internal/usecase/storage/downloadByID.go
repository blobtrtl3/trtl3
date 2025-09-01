package storage

import (
	"fmt"
	"os"
)

func (bs *BlobStorage) DownloadByID(id string) (*[]byte, error) {
	blob, err := os.ReadFile(fmt.Sprintf("/tmp/blobs/%s", id))
	if err != nil {
		return nil, err
	}

	return &blob, nil
}
