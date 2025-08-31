package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindByBucket(bucket string) (*[]domain.BlobInfo, error) {
	rows, err := bs.db.Query("SELECT * FROM blobsinfo WHERE bucket=?", bucket)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bis []domain.BlobInfo

	for rows.Next() {
  	var bi domain.BlobInfo

  	if err := rows.Scan(&bi.ID, &bi.Bucket, &bi.Mime, &bi.Size, &bi.CreatedAt); err != nil {
			return nil, err
		}

		bis = append(bis, bi)
	}

	return &bis, nil
}
