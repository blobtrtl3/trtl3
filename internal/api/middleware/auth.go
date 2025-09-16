package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	secretToken := os.Getenv("TOKEN")
	if secretToken == "" {
		secretToken = "trtl3"
	}

	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")

		if !strings.HasPrefix(bearer, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid token"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(bearer, "Bearer ")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "authorization header missing or invalid"})
			c.Abort()
			return
		}

		if token != secretToken {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
