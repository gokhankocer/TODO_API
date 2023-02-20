package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func fromJSON(body []byte) gin.H {
	var data gin.H
	json.Unmarshal(body, &data)
	return data
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func MockKafkaProducer(topic string, msg interface{}) {

}

func TestSignup(t *testing.T) {
	recorder := httptest.NewRecorder()

	repo := new(MockUserRepository)
	repo.On("CreateUser", mock.Anything).Return(nil)
	c, _ := gin.CreateTestContext(recorder)
	c.Set("UserRepository", repo)
	c.Set("KafkaProducer", MockKafkaProducer)

	user := entities.User{
		Name:     "test_user",
		Password: "test_password",
		Email:    "test_email@example.com",
	}

	c.Request, _ = http.NewRequest(http.MethodPost, "/signup", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(toJSON(user)))
	handlers.Signup(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
	body, _ := ioutil.ReadAll(recorder.Body)

	var response map[string]interface{}
	json.Unmarshal(body, &response)
	expectedResponse := gin.H{"message": "User created successfully"}
	assert.Equal(t, expectedResponse, response)
	repo.AssertExpectations(t)
}
