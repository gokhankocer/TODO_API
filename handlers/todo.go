package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/models"
)

func AddTodo(c *gin.Context) {
	var requestTodo models.Todo

	c.BindJSON(&requestTodo)
	todo := entities.Todo{
		Status:      &requestTodo.Status,
		Description: &requestTodo.Description,
	}

	database.DB.Create(&todo)
	c.JSON(http.StatusOK, &todo)
}

func GetTodos(c *gin.Context) {
	todos := []models.Todo{}
	database.DB.Find(&todos)
	c.JSON(200, todos)
}
