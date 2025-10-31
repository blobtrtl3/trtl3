package blob
// TODO: do better error messages

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/infra/cache"
	"github.com/blobtrtl3/trtl3/internal/shared"
	"github.com/blobtrtl3/trtl3/pkg/domain"
)

type Service interface {
	Save(bucket, mime string, size int64, r io.Reader) (*domain.BlobInfo, error)
	FindByBucket(bucket string) ([]domain.BlobInfo, error)
	FindUnique(bucket, id string) (*domain.BlobInfo, error)
	Delete(bucket, id string) (bool, error)
	Download(bucket, id string) ([]byte, error)
	Serve(bucket, id string) (*ServeInfo, error)
	Sign(bucket, id string, TTL int, once bool) (string, error)
}

type service struct {
	repo *Repository
	signaturesCache cache.SignaturesCache
	queue       *Queue
}

func NewService(
	r *Repository,
	sc cache.SignaturesCache,
	q *Queue,
) Service {
	return &service{
		repo: r,
		signaturesCache: sc,
		queue:       q,
	}
}

func (s *service) Save(bucket, mime string, size int64, r io.Reader) (*domain.BlobInfo, error) {
	blobInfo := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    bucket,
		Mime:      mime,
		CreatedAt: time.Now(),
		Size:      size, // NOTE: size in bytes value
	}

	if err := s.queue.Append(blobInfo, r); err != nil {
		return &domain.BlobInfo{}, err
	}

	return blobInfo, nil
}

func (s *service) FindByBucket(bucket string) ([]domain.BlobInfo, error) {
	if bucket == "" {
		return []domain.BlobInfo{}, fmt.Errorf("the bucket field sent is empty")
	}

	blobsInfos, err := s.repo.FindByBucket(bucket)
	if err != nil {
		return []domain.BlobInfo{}, err
	}

	if blobsInfos == nil {
		return []domain.BlobInfo{}, fmt.Errorf("there is no metadata from blobs in the bucket %s", bucket)
	}

	return blobsInfos, nil
}

func (s *service) FindUnique(bucket, id string) (*domain.BlobInfo, error) {
	if id == "" || bucket == "" {
		return &domain.BlobInfo{}, fmt.Errorf("the bucket or id field is empty %s", bucket)
	}

	blobInfo, err := s.repo.FindUnique(bucket, id)
	if err != nil {
		return &domain.BlobInfo{}, err
	}

	return blobInfo, nil
}

func (s *service) Download(bucket, id string) ([]byte, error) {
	if id == "" || bucket == "" {
		return []byte{}, fmt.Errorf("the bucket or id field is empty %s", bucket)
	}

	b, err := s.repo.Download(bucket, id)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (s *service) Delete(bucket, id string) (bool, error) {
	if id == "" || bucket == "" {
		return false, fmt.Errorf("the bucket or id field is empty %s", bucket)
	}

	_, err := s.repo.Delete(bucket, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

type ServeInfo struct {
	Key  string
	Mime string
	Path string
}

func (s *service) Serve(bucket, id string) (*ServeInfo, error) {
	if id == "" || bucket == "" {
		return &ServeInfo{}, fmt.Errorf("the bucket or id field is empty %s", bucket)
	}

	info, err := s.repo.FindUnique(bucket, id)
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

func (s *service) Sign(bucket, id string, TTL int, once bool) (string, error) {
	if _, err := s.repo.FindUnique(bucket, id); err != nil {
		return "", err
	}

	now := time.Now()
	signature := fmt.Sprintf("%s%s", shared.GenShortID(), now.Format("050204")) // format to SSDDMM

	s.signaturesCache.Set(
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
