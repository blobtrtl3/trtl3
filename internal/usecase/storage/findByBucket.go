package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindByBucket(bucket string) (*[]domain.BlobInfo, error) {
	var bi []domain.BlobInfo

	if err := bs.db.QueryRow("SELECT * FROM blobsinfo WHERE bucket=?)", bucket).Scan(&bi); err != nil {
		return nil, err
	}

	return &bi, nil
}
