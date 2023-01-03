package main

import (
	//"errors"
	"log"
	"os"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/gokhankocer/TODO-API/middleware"
	//"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func main() {

	database.ConnectPostgres()
	database.ConnectRedis()

	router := gin.Default()
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
	router.Use(sessions.Sessions("mysession", store))
	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", handlers.Signup)
	publicRoutes.POST("/login", handlers.Login)

	publicRoutes.GET("/logout", middleware.JWTAuthMiddleware(), handlers.Logout)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.GET("/todos", handlers.GetTodos)
	protectedRoutes.POST("todos", handlers.AddTodo)
	protectedRoutes.DELETE("todos/:id", handlers.DeleteTodo)
	protectedRoutes.GET("todos/:id", handlers.GetTodo)
	protectedRoutes.PATCH("todos/:id", handlers.UpdateTodo)
	protectedRoutes.GET("users/", handlers.GetUsers)
	protectedRoutes.GET("users/:id", handlers.GetUserById)
	protectedRoutes.PATCH("users/:id", handlers.UpdateUser)
	protectedRoutes.DELETE("users/:id", handlers.DeleteUser)

	//database.DB.Migrator().CreateTable(&entities.User{}, &entities.Todo{})
	//router.GET("/protected", handlers.TokenVerifiyMiddleWare(ProtectedEndpoint))//
	log.Fatal(router.Run("localhost:8080"))
}
