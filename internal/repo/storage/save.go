package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/shared"
)

func (bs *BlobStorage) Save(blobInfo *domain.BlobInfo, r io.Reader) (bool, error) {
	var exists bool

	if err := bs.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=? AND bucket=?)",
		blobInfo.ID, blobInfo.Bucket,
	).Scan(&exists); err != nil {
		return false, err
	}

	if exists {
		blobInfo.ID = shared.GenShortID()
		return bs.Save(blobInfo, r)
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

	out, err := os.Create(filepath.Join(bs.dir, shared.GenBlobName(blobInfo.Bucket, blobInfo.ID)))
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
