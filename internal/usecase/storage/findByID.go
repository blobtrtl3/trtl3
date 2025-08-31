package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindByID(id string) (*domain.BlobInfo, error) {
	var bi domain.BlobInfo

	if err := bs.db.QueryRow("SELECT 1 FROM blobsinfo WHERE id=?)", id).Scan(&bi); err != nil {
		return nil, err
	}

	return &bi, nil
}
