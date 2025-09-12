package jobs

import "time"

func (j *Jobs) cleanSignatures() {
	for _, key := range j.signatures.FindAll() {
		if j.signatures.Get(key).TTL.Compare(time.Now()) <= 0 {
			j.signatures.Delete(key)
		}
	}
}
