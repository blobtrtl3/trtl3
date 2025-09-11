package middleware

import (
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/gin-gonic/gin"
)

func SignMiddleware(hashmap map[string]domain.Signature) gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Query("sign")

		signature, exists := hashmap[sign]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid signature",
			})

			c.Abort()
			return
		}

		if signature.TTL.Compare(time.Now()) <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "expired signature",
			})
			c.Abort()
			return
		}

		if signature.Once == true {
			delete(hashmap, sign)
		}

		c.Set("bucket", signature.Bucket)
		c.Set("id", signature.ID)

		c.Next()
	}
}
