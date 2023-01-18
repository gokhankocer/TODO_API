package handlers

/*import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gokhankocer/TODO-API/models"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	newUser := models.UserRequest{
		Name:     "Goki",
		Email:    "goki@test.com",
		Password: "test",
	}
	writer := makeRequest("POST", "/auth/signup", newUser, false)
	assert.Equal(t, http.StatusOK, writer.Code)
}

func TestLogin(t *testing.T) {
	user := models.UserRequest{
		Name:     "Goki",
		Email:    "goki@test.com",
		Password: "test",
	}
	writer := makeRequest("POST", "/auth/login", user, true)
	assert.Equal(t, http.StatusOK, writer.Code)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	_, exists := response["jwt"]

	assert.Equal(t, false, exists)
}

func TestGetUsers(t *testing.T) {
	writer := makeRequest("GET", "/api/users", nil, true)
	assert.Equal(t, http.StatusOK, writer.Code)
}*/
