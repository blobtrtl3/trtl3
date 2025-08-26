package domain

import "time"

type BlobInfo struct {
	ID        string    `json:"id"`
	Bucket		string 		`json:"bucket"`
	Mime 			int       `json:"mime_type"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
