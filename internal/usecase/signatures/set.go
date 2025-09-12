package signatures

import "github.com/blobtrtl3/trtl3/internal/domain"

func (ms *MapSignatures) Set(key string, val domain.Signature) {
	ms.hm[key] = val
}
