package storage

import (
	"os"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

func (bs *BlobStorage) Save(bi *domain.BlobInfo, blob *[]byte) (bool, error) {
	var existID bool
	if err := bs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=?)", bi.ID).Scan(&existID); err != nil {
		return false, err
	}

	_, err := bs.db.Exec(
		"INSERT INTO blobsinfo VALUES(?, ?, ?, ?, ?)",
		bi.ID,
		bi.Bucket,
		bi.Mime,
		bi.Size,
		bi.CreatedAt,
	)
	if err != nil {
		return false, err
	}

	if err := os.WriteFile("/blob", *blob, 0644); err != nil {
		return false, err
	}

	return true, nil
}

