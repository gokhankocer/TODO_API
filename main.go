package main

import (
	//"errors"
	"log"
	//"net/http"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

type Todo struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func main() {
	router := gin.Default()
	//router.GET("/todos", getTodos)
	//router.POST("todos", addTodo)
	//router.GET("todos/:id", getTodo)
	//router.PATCH("todos/:id", updateTodoStatus)
	log.Fatal(router.Run("localhost:9090"))
}
