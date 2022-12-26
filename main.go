package main

import (
	//"errors"
	"log"
	//"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/handlers"
	//"github.com/gokhankocer/TODO-API/helper"
	//"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func main() {

	database.Connect()
	router := gin.Default()
	//publicRoutes := router.Group("/auth")
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)
	//protectedRoutes := router.Group("/api")
	//protectedRoutes.Use(JWTAuthMiddleware())
	router.GET("/todos", handlers.GetTodos)
	router.POST("todos", handlers.AddTodo)
	router.DELETE("todos/:id", handlers.DeleteTodo)
	router.GET("todos/:id", handlers.GetTodo)
	router.PATCH("todos/:id", handlers.UpdateTodo)

	//database.DB.Migrator().CreateTable(&entities.User{}, &entities.Todo{})
	//router.GET("/protected", handlers.TokenVerifiyMiddleWare(ProtectedEndpoint))//
	log.Fatal(router.Run("localhost:8080"))
}

/*func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.ValidateJWT(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Required"})
			c.Abort()
			return
		}
		c.Next()
	}
}*/
