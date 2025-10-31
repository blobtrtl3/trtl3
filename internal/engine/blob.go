package engine

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/shared"
)

type BlobEngine struct {
	db  *sql.DB
	dir string
}

func NewBlobEngine(db *sql.DB, dir string) *BlobEngine {
	return &BlobEngine{
		db:  db,
		dir: dir,
	}
}

func (be *BlobEngine) Save(blobInfo *domain.BlobInfo, r io.Reader) (bool, error) {
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

func (be *BlobEngine) FindAll() ([]domain.BlobInfo, error) {
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

func (be *BlobEngine) FindByBucket(bucket string) ([]domain.BlobInfo, error) {
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

func (be *BlobEngine) FindUnique(bucket string, id string) (*domain.BlobInfo, error) {
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

func (be *BlobEngine) Download(bucket string, id string) ([]byte, error) {
	blob, err := os.ReadFile(filepath.Join(be.dir, shared.GenBlobName(bucket, id)))
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (be *BlobEngine) Delete(bucket string, id string) (bool, error) {
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
