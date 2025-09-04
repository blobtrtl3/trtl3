package storage

import "github.com/blobtrtl3/trtl3/internal/domain"

func (bs *BlobStorage) FindAll() ([]domain.BlobInfo, error) {
	rows, err := bs.db.Query("SELECT * FROM blobsinfo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blobsInfos []domain.BlobInfo

	for rows.Next() {
		var blobInfo domain.BlobInfo

		if err := rows.Scan(&blobInfo.ID, &blobInfo.Bucket, &blobInfo.Mime, &blobInfo.Size, &blobInfo.CreatedAt); err != nil {
			return nil, err
		}

		blobsInfos = append(blobsInfos, blobInfo)
	}

	return blobsInfos, nil
}
