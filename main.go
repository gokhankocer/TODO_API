package main

import (
	//"errors"
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/handlers"
	//"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func main() {
	router := gin.Default()
	database.Connect()
	router.GET("/todos", handlers.GetTodos)
	router.POST("todos", handlers.AddTodo)
	router.DELETE("todos/:id", handlers.DeleteTodo)
	router.GET("todos/:id", handlers.GetTodo)
	router.PATCH("todos/:id", handlers.UpdateTodo)
	router.POST("/signup", handlers.Signup)
	database.DB.AutoMigrate(&entities.User{})
	router.POST("/login", handlers.Login)
	//router.GET("/protected", handlers.TokenVerifiyMiddleWare(ProtectedEndpoint))//
	log.Fatal(router.Run("localhost:8080"))
}
