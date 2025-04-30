package middleware

import (
	"github.com/gin-gonic/gin"
	"sky_storage_golang/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// try to get the cookie
		tokenStr, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication token required"})
			return
		}

		// validate token
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}
