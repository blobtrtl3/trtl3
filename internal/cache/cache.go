package cache

import "github.com/blobtrtl3/trtl3/internal/domain"

type SignaturesCache interface {
	Set(key string, val domain.Signature)
	Get(key string) *domain.Signature
	Delete(key string)
	FindAll() []string
}

type MemSignaturesCache map[string]domain.Signature

func NewMemSignaturesCache() SignaturesCache {
	return MemSignaturesCache{}
}

func (ms MemSignaturesCache) Set(key string, val domain.Signature) {
	ms[key] = val
}

func (ms MemSignaturesCache) Get(key string) *domain.Signature {
	if sig, ok := ms[key]; ok {
		return &sig
	}
	return nil
}

func (ms MemSignaturesCache) Delete(key string) {
	delete(ms, key)
}

func (ms MemSignaturesCache) FindAll() []string {
	keys := make([]string, 0, len(ms))
	for key := range ms {
		keys = append(keys, key)
	}
	return keys
}

