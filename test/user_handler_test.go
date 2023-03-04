package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/gokhankocer/TODO-API/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {

	userRepositoryMock := &mocks.UserRepositoryInterface{}
	userRepositoryMock.On("CreateUser", &entities.User{
		Name:     "Gokhan",
		Email:    "gokhan@test.com",
		Password: "123",
	}).Return(nil)

	userHandler := handlers.NewUserHandler(userRepositoryMock)
	router := gin.Default()
	router.POST("/auth/signup", userHandler.Signup)
	user := entities.User{
		Name:     "Gokhan",
		Email:    "gokhan@test.com",
		Password: "123",
	}
	requestBody, _ := json.Marshal(user)
	request, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	userRepositoryMock.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {

	userRepositoryMock := &mocks.UserRepositoryInterface{}

	var users = []entities.User{
		{
			Name:     "Gokhan",
			Email:    "gokhan@test.com",
			Password: "123",
		},
	}

	userRepositoryMock.On("GetUsers").Return(users, nil)

	userHandler := handlers.NewUserHandler(userRepositoryMock)
	router := gin.Default()
	router.GET("/api/users", userHandler.GetUsers)

	request, _ := http.NewRequest("GET", "/api/users", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	userRepositoryMock.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {

	userRepositoryMock := &mocks.UserRepositoryInterface{}

	user := entities.User{
		Name:     "Gokhan",
		Email:    "gokhan@test.com",
		Password: "123",
	}

	userRepositoryMock.On("GetUserByID", uint(1)).Return(&user, nil)

	userHandler := handlers.NewUserHandler(userRepositoryMock)
	router := gin.Default()

	router.GET("/api/users/:id", userHandler.GetUserById)

	request, _ := http.NewRequest("GET", "/api/users/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	userRepositoryMock.AssertExpectations(t)
}
