package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"
	"github.com/gokhankocer/TODO-API/models"
)

const userkey = "user"

func Signup(c *gin.Context) {
	var user entities.User
	err := c.Bind(&user)

	if err != nil {
		fmt.Println(err)
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

	err = user.VerifyPassword(body.Password)
	if err != nil {
		c.Abort()
	}

	jwt, err := helper.GenerateJwt(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Create Token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwt})

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	fmt.Println(user)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func GetUsers(c *gin.Context) {
	var users []entities.User

	val, err := database.RDB.Get(c, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	err := database.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	c.JSON(http.StatusOK, &users)
}

func GetUserById(c *gin.Context) {
	var user entities.User
	if err := database.DB.Where("id=?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not Found",
		})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func UpdateUser(c *gin.Context) { //password hashlenmiyor
	var user entities.User
	if err := database.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User Not Found"})
		return
	}
	currentUser, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	if currentUser.ID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	c.BindJSON(&user)
	database.DB.Save(&user)
	response := models.UserResponse{
		ID:    uint64(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	var user entities.User
	if err := database.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User Not Found"})
		return
	}
	currentUser, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	if currentUser.ID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	database.DB.Delete(&user)
	c.JSON(200, &user)
}
