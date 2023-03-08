package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/helper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		secret := os.Getenv("SECRET")
		url := os.Getenv("AUTHORIZE_URL")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK || token != secret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func SetUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := helper.CurrentUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
