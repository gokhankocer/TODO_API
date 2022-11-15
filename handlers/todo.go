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

func DeleteTodo(c *gin.Context) {
	var todo entities.Todo
	database.DB.Where("id = ?", c.Param("id")).Delete(&todo)
	c.JSON(200, &todo)
}

func UpdateTodo(c *gin.Context) {
	var todo entities.Todo
	database.DB.Where("id = ?", c.Param("id")).First(&todo)
	c.BindJSON(&todo)
	database.DB.Save(&todo)
	c.JSON(200, todo)
}

func GetTodo(c *gin.Context) {
	var todo entities.Todo
	database.DB.Where(&todo, "id = ?", c.Param("id"))
	c.JSON(200, &todo)
}

func GetTodos(c *gin.Context) {
	todos := []entities.Todo{}
	database.DB.Find(&todos)
	c.JSON(200, &todos)
}
