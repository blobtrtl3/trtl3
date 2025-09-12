package signatures

import (
	"github.com/blobtrtl3/trtl3/internal/domain"
)

type Signatures interface {
	Set(key string, val domain.Signature)
	Get(key string) *domain.Signature
	Delete(key string)
	FindAll() []domain.Signature
}

type MapSignatures struct {
	hm map[string]domain.Signature
}

func NewMapSignatures(hm map[string]domain.Signature) Signatures {
	return &MapSignatures{
		hm:  hm,
	}
}
