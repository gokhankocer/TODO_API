package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/mocks"
	"github.com/gokhankocer/TODO-API/routers"
)

func TestSignup(t *testing.T) {
	router := routers.SetupRouter()
	repoMock := mocks.NewRepositoryInterface(t)
	repoMock.On("CreateUser", &entities.User{
		Name:     "Gokhan",
		Email:    "gokhanel@test.com",
		Password: "1234",
	}).Return(nil)
	// repoMock.On("HashPassword", "1234").Return("hashed_password", nil)

	// Set the repository to the mock repository

	user := entities.User{
		Name:     "Gokhan",
		Email:    "gokhanel@test.com",
		Password: "1234",
	}
	requestBody, _ := json.Marshal(user)
	request, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	bodyString := string(body)
	fmt.Println(bodyString)

	assert.Equal(t, http.StatusOK, response.Code)
	repoMock.AssertExpectations(t)
}
