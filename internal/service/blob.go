// TODO: do better error messages
package service

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/engine"
	"github.com/blobtrtl3/trtl3/internal/infra/cache"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/internal/shared"
)

type BlobService interface {
	Save(bucket, mime string, size int64, r io.Reader) (bool, error)
	FindByBucket(bucket string) ([]domain.BlobInfo, error)
	FindUnique(bucket, id string) (*domain.BlobInfo, error)
	Delete(bucket, id string) (bool, error)
	Download(bucket, id string) ([]byte, error)
	Serve(bucket, id string) (*ServeInfo, error)
	Sign(bucket, id string, TTL int, once bool) (string, error)
}

type BlobServiceImpl struct {
	blobEngine      engine.BlobEngine
	signaturesCache cache.SignaturesCache
	blobQueue       queue.BlobQueue
}

func NewBlobServiceImpl(
	be engine.BlobEngine,
	sc cache.SignaturesCache,
	bq queue.BlobQueue,
) BlobService {
	return &BlobServiceImpl{
		blobEngine:      be,
		signaturesCache: sc,
		blobQueue:       bq,
	}
}

func (bs *BlobServiceImpl) Save(bucket, mime string, size int64, r io.Reader) (bool, error) {
	blobInfo := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    bucket,
		Mime:      mime,
		CreatedAt: time.Now(),
		Size:      size, // NOTE: size in bytes value
	}

	if err := bs.blobQueue.Append(blobInfo, r); err != nil {
		return false, err
	}

	return true, nil
}

func (bs *BlobServiceImpl) FindByBucket(bucket string) ([]domain.BlobInfo, error) {
	if bucket == "" {
		return []domain.BlobInfo{}, fmt.Errorf("the bucket field sent is empty")
	}

	blobsInfos, err := bs.blobEngine.FindByBucket(bucket)
	if err != nil {
		return []domain.BlobInfo{}, err
	}

	if blobsInfos == nil {
		return []domain.BlobInfo{}, fmt.Errorf("there is no metadata from blobs in the bucket %s", bucket)
	}

	return blobsInfos, nil
}

func (bs *BlobServiceImpl) FindUnique(bucket, id string) (*domain.BlobInfo, error) {
	if id == "" || bucket == "" {
		return &domain.BlobInfo{}, fmt.Errorf("the bucket or id field is empty", bucket)
	}

	blobInfo, err := bs.blobEngine.FindUnique(bucket, id)
	if err != nil {
		return &domain.BlobInfo{}, err
	}

	return blobInfo, nil
}

func (bs *BlobServiceImpl) Download(bucket, id string) ([]byte, error) {
	if id == "" || bucket == "" {
		return []byte{}, fmt.Errorf("the bucket or id field is empty", bucket)
	}

	b, err := bs.blobEngine.Download(bucket, id)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (bs *BlobServiceImpl) Delete(bucket, id string) (bool, error) {
	if id == "" || bucket == "" {
		return false, fmt.Errorf("the bucket or id field is empty", bucket)
	}

	_, err := bs.blobEngine.Delete(bucket, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

type ServeInfo struct {
	key  string
	mime string
	path string
}

func (bs *BlobServiceImpl) Serve(bucket, id string) (*ServeInfo, error) {
	if id == "" || bucket == "" {
		return &ServeInfo{}, fmt.Errorf("the bucket or id field is empty", bucket)
	}

	info, err := bs.blobEngine.FindUnique(bucket, id)
	if err != nil {
		return &ServeInfo{}, err
	}

	key := shared.GenBlobName(info.Bucket, info.ID)

	return &ServeInfo{
		key,
		info.Mime,
		filepath.Join("blobs", key),
	}, nil
}

func (bs *BlobServiceImpl) Sign(bucket, id string, TTL int, once bool) (string, error) {
	if _, err := bs.blobEngine.FindUnique(bucket, id); err != nil {
		return "", err
	}

	now := time.Now()
	signature := fmt.Sprintf("%s%s", shared.GenShortID(), now.Format("050204")) // format to SSDDMM

	bs.signaturesCache.Set(
		signature,
		domain.Signature{
			Bucket: bucket,
			ID:     id,
			TTL:    now.Add(time.Duration(TTL) * time.Minute),
			Once:   once,
		},
	)

	return signature, nil
}
