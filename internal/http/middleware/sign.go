package middleware

import (
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/infra"
	"github.com/gin-gonic/gin"
)

func SignMiddleware(s infra.SignaturesCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Query("sign")

		signature := s.Get(sign)
		if signature == nil {
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

		if signature.Once {
			s.Delete(sign)
		}

		c.Set("bucket", signature.Bucket)
		c.Set("id", signature.ID)

		c.Next()
	}
}
