package signatures

import "github.com/blobtrtl3/trtl3/internal/domain"

func (ms *MapSignatures) FindAll() []domain.Signature {
	data := make([]domain.Signature, 0, len(ms.hm))
	for _, sig := range ms.hm {
		data = append(data, sig)
	}
	return data
}

