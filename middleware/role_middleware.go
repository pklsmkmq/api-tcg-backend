package middleware

import "github.com/gin-gonic/gin"

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(403, gin.H{"message": "Admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
