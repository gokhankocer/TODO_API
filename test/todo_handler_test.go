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
	"github.com/gokhankocer/TODO-API/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddTodo(t *testing.T) {
	todoRepositoryMock := &mocks.TodoRepositoryInterface{}
	userRepositoryMock := &mocks.UserRepositoryInterface{}

	user := entities.User{

		Name:     "Gokhan ",
		Email:    "test@gmail.com",
		Password: "password",
		IsActive: true,
	}

	todo := entities.Todo{
		ID:          1,
		Description: "Test Todo",
		Status:      "Todo",
		UserID:      user.ID,
	}

	requestTodo := models.PostTodoRequest{
		Description: "Test Todo",
		Status:      "Todo",
	}

	userRepositoryMock.On("GetUserByID", mock.Anything).Return(&user, nil)
	todoRepositoryMock.On("AddTodo", todo).Return(&todo, nil)

	handler := handlers.NewTodoHandler(todoRepositoryMock, userRepositoryMock)
	router := gin.Default()
	router.POST("/api/todos", handler.AddTodo)

	payload, _ := json.Marshal(requestTodo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(payload))

	// Serve the request and get the response
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Assert the response code is 201 (created)

	// Assert the mock expectations
	todoRepositoryMock.AssertExpectations(t)
	userRepositoryMock.AssertExpectations(t)
}
func TestGetTodos(t *testing.T) {
	todoRepositoryMock := &mocks.TodoRepositoryInterface{}
	userRepositoryMock := &mocks.UserRepositoryInterface{}

	user := entities.User{

		Name:     "Gokhan ",
		Email:    "test@gmail.com",
		Password: "password",
		IsActive: true,
	}

	todo := entities.Todo{
		ID:          1,
		Description: "Test Todo",
		Status:      "Todo",
		UserID:      user.ID,
	}

	userRepositoryMock.On("GetUserByID", mock.Anything).Return(&user, nil)
	todoRepositoryMock.On("GetTodos", todo).Return(todo, nil)

	handler := handlers.NewTodoHandler(todoRepositoryMock, userRepositoryMock)
	router := gin.Default()

	router.GET("/api/todos", handler.GetTodos)

	request, _ := http.NewRequest("GET", "/api/todos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	// Assert the expectations of the mock repository
	todoRepositoryMock.AssertExpectations(t)
}
