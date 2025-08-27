package storage

import (
	"fmt"
	"os"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

func (bs *BlobStorage) Save(bi *domain.BlobInfo, blob *[]byte) (bool, error) {
	var existID bool
	if err := bs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=?)", bi.ID).Scan(&existID); err != nil {
		return false, err
	}

  if existID {
    return false, fmt.Errorf("blob with id %s already exists", bi.ID)
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

	if err := os.WriteFile(
		fmt.Sprintf("/tmp/blobs/%s", bi.ID),
		*blob,
		os.ModePerm,
	); err != nil {
		return false, err
	}

	return true, nil
}

