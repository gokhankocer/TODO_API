package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
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
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
