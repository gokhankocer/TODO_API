package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"
	"github.com/gokhankocer/TODO-API/services/email"
	"github.com/gokhankocer/TODO-API/services/kafka_client"

	"github.com/gokhankocer/TODO-API/models"
	"github.com/gokhankocer/TODO-API/repository"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserRepository repository.UserRepositoryInterface
}

func NewUserHandler(userRepository repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		UserRepository: userRepository,
	}
}

func (handler *UserHandler) Signup(c *gin.Context) {
	var user entities.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	if err := handler.UserRepository.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	go kafka_client.Producer("mail", user)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (handler *UserHandler) GetUsers(c *gin.Context) {
	users, err := handler.UserRepository.GetUsers()
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Request"})
		return
	}

	c.JSON(http.StatusOK, &users)
}

func (handler *UserHandler) GetUserById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	user, err := handler.UserRepository.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not Found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (handler *UserHandler) UpdateUser(c *gin.Context) {
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
	user, err := handler.UserRepository.GetUserByID(uintID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User Not Found"})
		return
	}
	currentUserID, _ := helper.CurrentUser(c)
	currentUser, err := handler.UserRepository.FindUserById(currentUserID)

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

	if err := handler.UserRepository.UpdateUser(uintID, user); err != nil {
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

func (handler *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	currentUserID, _ := helper.CurrentUser(c)
	currentUser, err := handler.UserRepository.FindUserById(currentUserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}

	if uint(currentUser.ID) != uint(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}

	if err := handler.UserRepository.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
func (handler *UserHandler) ResetPassword(c *gin.Context) {
	var request models.UserRequest
	if c.BindJSON(&request) != nil {
		log.Println("error", "Body Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	user, err := handler.UserRepository.FindUserByEmail(request.Email)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	resetPasswordToken := uuid.New().String()
	user.ResetPasswordToken = resetPasswordToken
	if err := handler.UserRepository.UpdateUser(user.ID, &user); err != nil {
		log.Println("error", "Failed to update reset password token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reset password token"})
		return
	}

	resetPasswordLink := fmt.Sprintf("http://localhost:3000/reset_password/%s", resetPasswordToken)

	// Send reset password email to the user's email address
	email.SendResetPasswordEmail(user.Email, resetPasswordLink)
	log.Println("kafka")
	c.JSON(http.StatusOK, gin.H{"message": "Reset password email sent"})
}

func (handler *UserHandler) ConfirmResetPassword(c *gin.Context) {
	token := c.Param("token")

	user, err := handler.UserRepository.FindUserByResetPasswordToken(token)
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
	if err := handler.UserRepository.UpdateUser(user.ID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Write the code to update the user's password

func (handler *UserHandler) Activate(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User Id"})
		return
	}
	user, err := handler.UserRepository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user.IsActive = true
	if err := handler.UserRepository.UpdateUser(uint(userID), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully activated"})

}
