package domain

import "time"

type Signature struct {
	Bucket string
	ID     string
	TTL    time.Time
	Once   bool // if once is true after one access should delete the signature
}
