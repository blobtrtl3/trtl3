package domain

import "time"

type Signature struct {
	TTL time.Time
	Once bool // if once is true after one access should delete the signature
}
