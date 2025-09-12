package signatures

import "github.com/blobtrtl3/trtl3/internal/domain"

func (ms *MapSignatures) Get(key string) (domain.Signature) {
	return ms.hm[key]
}
