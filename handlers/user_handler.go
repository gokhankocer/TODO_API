package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	user := entities.User{Name: body.Name, Email: body.Email, Password: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user entities.User
	database.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed To Create Token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}
