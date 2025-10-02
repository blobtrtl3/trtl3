package blob

import (
	"io"
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/shared"
)

func (be *BlobEngineImpl) Save(blobInfo *domain.BlobInfo, r io.Reader) (bool, error) {
	var exists bool

	if err := be.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM blobsinfo WHERE id=? AND bucket=?)",
		blobInfo.ID, blobInfo.Bucket,
	).Scan(&exists); err != nil {
		return false, err
	}

	if exists {
		blobInfo.ID = shared.GenShortID()
		return be.Save(blobInfo, r)
	}

	tx, err := be.db.Begin()
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

	out, err := os.Create(filepath.Join(be.dir, shared.GenBlobName(blobInfo.Bucket, blobInfo.ID)))
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

func (be *BlobEngineImpl) FindAll() ([]domain.BlobInfo, error) {
	rows, err := be.db.Query("SELECT * FROM blobsinfo")
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

func (be *BlobEngineImpl) FindByBucket(bucket string) ([]domain.BlobInfo, error) {
	rows, err := be.db.Query("SELECT * FROM blobsinfo WHERE bucket=?", bucket)
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

func (be *BlobEngineImpl) FindUnique(bucket string, id string) (*domain.BlobInfo, error) {
	var blobInfo domain.BlobInfo

	if err := be.db.QueryRow("SELECT * FROM blobsinfo WHERE bucket=? AND id=?", bucket, id).Scan(
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

func (be *BlobEngineImpl) Download(bucket string, id string) ([]byte, error) {
	blob, err := os.ReadFile(filepath.Join(be.dir, shared.GenBlobName(bucket, id)))
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (be *BlobEngineImpl) Delete(bucket string, id string) (bool, error) {
	tx, err := be.db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(
		"DELETE FROM blobsinfo WHERE bucket=? AND id=?",
		bucket,
		id,
	)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := os.Remove(filepath.Join(be.dir, shared.GenBlobName(bucket, id))); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
