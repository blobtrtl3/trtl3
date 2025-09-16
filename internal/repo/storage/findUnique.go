package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindUnique(bucket string, id string) (*domain.BlobInfo, error) {
	var blobInfo domain.BlobInfo

	if err := bs.db.QueryRow("SELECT * FROM blobsinfo WHERE bucket=? AND id=?", bucket, id).Scan(
		&blobInfo.ID,
		&blobInfo.Bucket,
		&blobInfo.Mime,
		&blobInfo.Size,
		&blobInfo.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &blobInfo, nil
}
