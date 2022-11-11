package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/models"
)

func AddTodo(c *gin.Context) {
	var todo models.Todo
	c.BindJSON(&todo)
	database.DB.Create(&todo)
	c.JSON(http.StatusOK, &todo)
}
