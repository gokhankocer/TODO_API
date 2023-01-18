package handlers

/*import (
	"os"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/models"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func router() *gin.Engine {
	router := gin.Default()
	publicRoutes := router.Group("auth")
	publicRoutes.POST("/signup", Signup)
	publicRoutes.POST("/login", Login)

	protectedRoutes := router.Group("/api")
	//protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("todos", AddTodo)
	protectedRoutes.GET("/todos", GetTodos)
	protectedRoutes.GET("/users", GetUsers)


	return router
}

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

}

func teardown() {
	migrator := database.DB.Migrator()
	migrator.DropTable(&entities.User{})
	migrator.DropTable(&entities.Todo{})
}

func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Autherization", "Bearer"+bearerToken())
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

// mockery
func bearerToken() string {
	user := models.UserRequest{
		Name:     "Goki",
		Email:    "goki@test.com",
		Password: "test",
	}
	writer := makeRequest("POST", "/auth/login", user, false)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["jwt"]
}*/
