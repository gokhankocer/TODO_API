package middlerware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/helper"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.ValidateJWT(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
