package storage

import (
	"fmt"
	"os"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

func (bs *BlobStorage) Save(blobInfo *domain.BlobInfo, blobBytes *[]byte) (bool, error) {
	var existID bool

	if err := bs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=?)", blobInfo.ID).Scan(&existID); err != nil {
		return false, err
	}

	if existID {
		return false, fmt.Errorf("blob with id %s already exists", blobInfo.ID)
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
		*blobBytes,
		os.ModePerm,
	); err != nil {
		return false, err
	}

	return true, nil
}
