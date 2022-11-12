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
	//router.GET("todos/:id", getTodo)
	//router.PATCH("todos/:id", updateTodoStatus)
	log.Fatal(router.Run("localhost:8080"))
}
