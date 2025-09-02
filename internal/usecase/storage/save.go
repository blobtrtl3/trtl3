package storage

import (
	"fmt"
	"os"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/shared"
)

func (bs *BlobStorage) Save(blobInfo *domain.BlobInfo, blobBytes []byte) (bool, error) {
	var exists bool

	if err := bs.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=? AND bucket=?)",
		blobInfo.ID, blobInfo.Bucket,
	).Scan(&exists); err != nil {
		return false, err
	}

	if exists {
		blobInfo.ID = shared.GenShortID()
		return bs.Save(blobInfo, blobBytes)
	}

	_, err := bs.db.Exec(
		"INSERT INTO blobsinfo VALUES(?, ?, ?, ?, ?)",
		blobInfo.ID,
		blobInfo.Bucket,
		blobInfo.Mime,
		blobInfo.Size,
		blobInfo.CreatedAt,
	)
	if err != nil {
		return false, err
	}

	if err := os.WriteFile(
		fmt.Sprintf("/tmp/blobs/%s", blobInfo.ID),
		blobBytes,
		os.ModePerm,
	); err != nil {
		return false, err
	}

	return true, nil
}
