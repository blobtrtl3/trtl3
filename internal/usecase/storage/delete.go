package storage

import (
	"fmt"
	"os"
)

func (bs *BlobStorage) Delete(bucket string, id string) (bool, error) {
	_, err := bs.db.Exec(
		"DELETE FROM blobsinfo WHERE bucket=? AND id=?",
		bucket,
		id,
	)
	if err != nil {
		return false, err
	}

	if err := os.Remove(fmt.Sprintf("/tmp/blobs/%s", id)); err != nil {
		return false, err
	}

	return true, nil
}
