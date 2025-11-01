package blob

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"

	"github.com/blobtrtl3/trtl3/internal/shared"
	"github.com/blobtrtl3/trtl3/pkg/domain"
)

type Repository struct {
	db  *sql.DB
	dir string
}

func NewRepository(db *sql.DB, dir string) *Repository {
	return &Repository{
		db:  db,
		dir: dir,
	}
}

func (r *Repository) Save(blobInfo *domain.BlobInfo, ior io.Reader) (bool, error) {
	tx, err := r.db.Begin()
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

	out, err := os.Create(filepath.Join(r.dir, shared.GenBlobName(blobInfo.Bucket, blobInfo.ID)))
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer out.Close()

	if _, err := io.Copy(out, ior); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) FindAll() ([]domain.BlobInfo, error) {
	rows, err := r.db.Query("SELECT * FROM blobsinfo")
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

func (r *Repository) FindByBucket(bucket string) ([]domain.BlobInfo, error) {
	rows, err := r.db.Query("SELECT * FROM blobsinfo WHERE bucket=?", bucket)
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

func (r *Repository) FindUnique(bucket string, id string) (*domain.BlobInfo, error) {
	var blobInfo domain.BlobInfo

	if err := r.db.QueryRow("SELECT * FROM blobsinfo WHERE bucket=? AND id=?", bucket, id).Scan(
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

func (r *Repository) Download(bucket string, id string) ([]byte, error) {
	blob, err := os.ReadFile(filepath.Join(r.dir, shared.GenBlobName(bucket, id)))
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (r *Repository) Delete(bucket string, id string) (bool, error) {
	tx, err := r.db.Begin()
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

	if err := os.Remove(filepath.Join(r.dir, shared.GenBlobName(bucket, id))); err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
