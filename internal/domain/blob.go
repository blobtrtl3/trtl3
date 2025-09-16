package domain

import "time"

type BlobInfo struct {
	ID        string    `json:"id"`
	Bucket    string    `json:"bucket"`
	Mime      string    `json:"mime_type"`
	Size      int64       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
