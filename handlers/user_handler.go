package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"
	"github.com/gokhankocer/TODO-API/kafka_service/kafka"
	"github.com/gokhankocer/TODO-API/models"
	"github.com/gokhankocer/TODO-API/repository"
	"github.com/google/uuid"
)

func Signup(c *gin.Context) {

	var user entities.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := repository.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	go kafka.Producer("new_user", user)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func GetUsers(c *gin.Context) {
	users, err := repository.GetUsers()
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Request"})
		return
	}

	c.JSON(http.StatusOK, &users)
}

func GetUserById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	user, err := repository.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not Found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	var userReq models.UpdateUserRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	uintID := uint(id)
	user, err := repository.GetUserByID(uintID)
	if err != nil {
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

	if userReq.Password != "" {
		user.HashPassword(userReq.Password)
	}
	if userReq.Name != "" {
		user.Name = userReq.Name
	}
	if userReq.Email != "" { // TODO check email already exists.
		user.Email = userReq.Email
	}

	if err := repository.UpdateUser(uintID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update Error"})
		return
	}
	response := models.UserResponse{
		ID:    uint64(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	currentUser, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}

	if uint(currentUser.ID) != uint(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}

	if err := repository.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
func ResetPassword(c *gin.Context) {
	var request models.UserRequest
	if c.BindJSON(&request) != nil {
		log.Println("error", "Body Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user, err := repository.FindUserByEmail(request.Email)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	resetPasswordToken := uuid.New().String()
	user.ResetPasswordToken = resetPasswordToken
	if err := repository.UpdateUser(user.ID, &user); err != nil {
		log.Println("error", "Failed to update reset password token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reset password token"})
		return
	}

	resetPasswordLink := fmt.Sprintf("http://localhost:3000/reset_password/%s", resetPasswordToken)

	// Send reset password email to the user's email address
	kafka.SendResetPasswordEmail(user.Email, resetPasswordLink)
	log.Println("kafka")
	c.JSON(http.StatusOK, gin.H{"message": "Reset password email sent"})
}

func ConfirmResetPassword(c *gin.Context) {
	token := c.Param("token")

	user, err := repository.FindUserByResetPasswordToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var body models.ConfirmResetPasswordRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user.HashPassword(body.Password)
	user.ResetPasswordToken = ""
	if err := repository.UpdateUser(user.ID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Write the code to update the user's password
