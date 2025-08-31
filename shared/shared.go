package shared

import (
	"crypto/rand"
	"encoding/base64"
)

func GenShortID() string {
	n := 9
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}
