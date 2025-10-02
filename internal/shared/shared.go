package shared

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenShortID() string {
	n := 6
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func GenBlobName(bucket string, id string) string {
	return fmt.Sprintf("%s_%s", bucket, id)
}
