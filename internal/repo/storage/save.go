package storage

import (
	"os"
	"path/filepath"

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

	tx, err := bs.db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(
		"INSERT INTO blobsinfo VALUES(?, ?, ?, ?, ?)",
		blobInfo.ID,
		blobInfo.Bucket,
		blobInfo.Mime,
		blobInfo.Size,
		blobInfo.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := os.WriteFile(
		filepath.Join(bs.dir, shared.GenBlobName(blobInfo.Bucket, blobInfo.ID)),
		blobBytes,
		os.ModePerm,
	); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
