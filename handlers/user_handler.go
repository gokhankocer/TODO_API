package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"
	"github.com/gokhankocer/TODO-API/models"
)

func Signup(c *gin.Context) {
	var user entities.User
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	var body models.UserRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user, err := entities.FindUserByName(body.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Name"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) //VerifyPassword u gormuyor
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}
	jwt, err := helper.GenerateJwt(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Create Token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwt})

}
