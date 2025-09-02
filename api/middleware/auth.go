package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var secretToken = "trtl3" // TODO: take token from env

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

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
