package domain

import "time"

// type BlobMeta struct { // TODO: rename
type BlobInfo struct {
	ID        string    `json:"id"`
	Bucket    string    `json:"bucket"`
	Mime      string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type Signature struct {
	Bucket string
	ID     string
	TTL    time.Time
	Once   bool // if once is true after one access should delete the signature
}
