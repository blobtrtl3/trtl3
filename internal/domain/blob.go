package domain

import "time"

type Blob struct {
	ID        string    `json:"id"`
	MimeType  int       `json:"mime_type"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
