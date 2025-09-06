package storage

import (
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/shared"
)

func (bs *BlobStorage) DownloadByID(bucket string, id string) ([]byte, error) {
	blob, err := os.ReadFile(filepath.Join(bs.dir, shared.GenBlobName(bucket, id)))
	if err != nil {
		return nil, err
	}

	return blob, nil
}
