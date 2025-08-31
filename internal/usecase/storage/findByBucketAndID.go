package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindByBucketAndID(bucket string, id string) (*domain.BlobInfo, error) {
	var bi domain.BlobInfo

	if err := bs.db.QueryRow("SELECT * FROM blobsinfo WHERE bucket=? AND id=?)", bucket, id).Scan(&bi); err != nil {
		return nil, err
	}

	return &bi, nil
}
