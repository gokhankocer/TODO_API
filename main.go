package main

import (
	//"errors"

	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/handlers"
	"github.com/gokhankocer/TODO-API/middleware"
	"github.com/gokhankocer/TODO-API/repository"
	"github.com/joho/godotenv"
	//"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func main() {

	loadEnv()
	database.ConnectPostgres()

	migrateFlag := flag.Bool("migrate", false, "migrate argument")
	flag.Parse()
	if *migrateFlag {
		fmt.Println("Migrate started!")
		database.DB.Migrator().DropTable(&entities.User{}, &entities.Todo{})
		database.DB.Migrator().CreateTable(&entities.User{}, &entities.Todo{})
		return
	}

	database.ConnectRedis()
	// go kafka.Consume(context.Background(), "email")
	// go kafka.Consume(context.Background(), "mail")
	userRepository := repository.NewUserRepository(database.DB)
	todoRepository := repository.NewTodoRepository(database.DB)
	userHandler := handlers.NewUserHandler(userRepository)
	todoHandler := handlers.NewTodoHandler(todoRepository, userRepository)

	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", userHandler.Signup)
	//publicRoutes.POST("/login", handlers.Login)

	//publicRoutes.GET("/logout", middleware.AuthMiddleware(), handlers.Logout)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware())
	protectedRoutes.GET("/todos", todoHandler.GetTodos)
	protectedRoutes.POST("todos", todoHandler.AddTodo)
	protectedRoutes.DELETE("todos/:id", todoHandler.DeleteTodo)
	protectedRoutes.GET("todos/:id", todoHandler.GetTodo)
	protectedRoutes.PATCH("todos/:id", todoHandler.UpdateTodo)
	protectedRoutes.GET("users/", userHandler.GetUsers)
	protectedRoutes.GET("users/:id", userHandler.GetUserById)
	protectedRoutes.PATCH("users/:id", userHandler.UpdateUser)
	protectedRoutes.DELETE("users/:id", userHandler.DeleteUser)

	// router.GET("/api/activate/:id", kafka.Activate)
	router.PATCH("/reset_password/:token", userHandler.ConfirmResetPassword)
	router.POST("/reset_password/", userHandler.ResetPassword)
	log.Fatal(router.Run("0.0.0.0:3000"))

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file")
	} else {
		log.Print("Env successfully loaded")
	}

}
