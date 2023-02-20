package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/gokhankocer/TODO-API/helper"
)

func TestGetTodos(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/todos", handlers.GetTodos)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos", nil)

	// Test with no errors
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but received %d", http.StatusOK, w.Code)
	}

	// Test with errors
	mockUserFunc := func(*gin.Context) (interface{}, error) {
		return nil, fmt.Errorf("Mock Error")
	}
	router.GET("/todos", GetTodosWithMockUserFunc(mockUserFunc))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/todos", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d but received %d", http.StatusInternalServerError, w.Code)
	}
}

func GetTodosWithMockUserFunc(mockUserFunc func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := mockUserFunc(c)
		if err != nil {
			log.Println("error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
func TestAddTodo(t *testing.T) {
	r := gin.Default()
	r.POST("/todos", handlers.AddTodo)
	var user entities.User
	token, err := helper.GenerateJwt(user)
	if err != nil {
		t.Errorf("Error generating JWT: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todos", strings.NewReader(`{"status": "pending", "description": "Write code"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected response code 201, but got %v", w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	if response["data"] == nil {
		t.Error("Response data is nil")
	}
}
