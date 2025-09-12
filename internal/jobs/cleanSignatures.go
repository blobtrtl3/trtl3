package jobs

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
)

func (j *Jobs) cleanSignatures(hashmap map[string]domain.Signature) {
  for key, val := range hashmap {
		if val.TTL.Compare(time.Now()) <= 0 {
			delete(hashmap, key)
		}
  }
}
