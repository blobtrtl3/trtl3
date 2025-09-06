package storage

import (
	"fmt"
	"os"
)

func (bs *BlobStorage) Delete(bucket string, id string) (bool, error) {
	tx, err := bs.db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(
		"DELETE FROM blobsinfo WHERE bucket=? AND id=?",
		bucket,
		id,
	)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := os.Remove(fmt.Sprintf("/tmp/blobs/%s_%s", bucket, id)); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
