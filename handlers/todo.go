package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/models"
)

func AddTodo(c *gin.Context) {
	requestTodo := &models.PostTodoRequest{}
	if err := c.BindJSON(requestTodo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid Request",
		})
		return
	}

	err := requestTodo.Validate(strfmt.Default)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing Fields",
		})
		return
	}
	todo := entities.Todo{
		Status:      &requestTodo.Status,
		Description: &requestTodo.Description,
	}

	database.DB.Create(&todo)
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      *todo.Status,
		Description: *todo.Description,
	}

	c.JSON(http.StatusOK, &responseTodo)

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
	database.DB.Where("id = ?", c.Param("id")).First(&todo)
	c.JSON(200, &todo)
}

func GetTodos(c *gin.Context) {
	todos := []entities.Todo{}
	database.DB.Find(&todos)
	c.JSON(200, &todos)
}
