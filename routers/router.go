package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/gokhankocer/TODO-API/kafka_service/kafka"
	"github.com/gokhankocer/TODO-API/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", handlers.Signup)
	//publicRoutes.POST("/login", handlers.Login)
	//publicRoutes.GET("/logout", middleware.AuthMiddleware(), handlers.Logout)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware())
	protectedRoutes.GET("/todos", handlers.GetTodos)
	protectedRoutes.POST("todos", handlers.AddTodo)
	protectedRoutes.DELETE("todos/:id", handlers.DeleteTodo)
	protectedRoutes.GET("todos/:id", handlers.GetTodo)
	protectedRoutes.PATCH("todos/:id", handlers.UpdateTodo)

	protectedRoutes.GET("users/", handlers.GetUsers)
	protectedRoutes.GET("users/:id", handlers.GetUserById)
	protectedRoutes.PATCH("users/:id", handlers.UpdateUser)
	protectedRoutes.DELETE("users/:id", handlers.DeleteUser)

	router.GET("/api/activate/:id", kafka.Activate)
	router.PATCH("/reset_password/:token", handlers.ConfirmResetPassword)
	router.POST("/reset_password/", handlers.ResetPassword)
	return router
}
