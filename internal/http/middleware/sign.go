package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SignMiddleware(ctx context.Context, r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Query("sign")

		res, err := r.Get(ctx, sign).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid signature",
			})

			c.Abort()
			return
		}

    var signature domain.Signature

		if err := json.Unmarshal([]byte(res), &signature); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
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
			if err := r.Del(ctx, sign).Err(); err != nil {
				// TODO: handler this
			}
		}

		c.Set("bucket", signature.Bucket)
		c.Set("id", signature.ID)

		c.Next()
	}
}
