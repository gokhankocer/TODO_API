package main

import (
	//"errors"
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/handlers"
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
	log.Fatal(router.Run("localhost:8080"))
}
