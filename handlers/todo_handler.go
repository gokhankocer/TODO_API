package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"

	"github.com/gokhankocer/TODO-API/models"
	"github.com/gokhankocer/TODO-API/repository"
)

/*type TodoHandler struct {
	TodoRepository repository.ToDoRepoInterface
}

func CreateHandeler(TodoRepo repository.ToDoRepoInterface) *TodoHandler {
	return &TodoHandler{TodoRepository: TodoRepo}
}*/

func AddTodo(c *gin.Context) {
	var requestTodo models.PostTodoRequest

	if err := c.BindJSON(&requestTodo); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid Request Payload",
		})

		return
	}

	err := requestTodo.Validate(strfmt.Default)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Missing Fields",
		})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		log.Println("error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Logged in"})
		return
	}
	var todo = entities.Todo{
		Status:      requestTodo.Status,
		Description: requestTodo.Description,
	}

	todo.UserID = user.ID
	savedTodo, err := repository.AddTodo(&todo)

	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": savedTodo})
}

func DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := repository.GetTodo(uint(id))
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Todo Not Found"})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	repository.DeleteTodo(todo)
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      todo.Status,
		Description: todo.Description,
	}
	c.JSON(http.StatusNoContent, &responseTodo)
}

func UpdateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := repository.GetTodo(uint(id))
	if err != nil {
		log.Println("error", err)

		c.JSON(http.StatusNotFound, gin.H{"message": "Todo Not Found"})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	c.BindJSON(&todo)
	repository.UpdateTodo(todo)
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      todo.Status,
		Description: todo.Description,
	}
	c.JSON(200, responseTodo)
}

func GetTodo(c *gin.Context) {
	user, err := helper.CurrentUser(c)
	if err != nil {
		log.Println("error getting current user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}
	log.Println("current user:", user)
	if user.ID == 0 {
		log.Println("current user is nil")
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("error", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Todo ID"})
		return
	}
	todo, err := repository.GetTodo(uint(id))
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo Not Found"})
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      todo.Status,
		Description: todo.Description,
	}
	c.JSON(http.StatusOK, &responseTodo)

}

func GetTodos(c *gin.Context) {
	user, err := helper.CurrentUser(c)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}

	c.JSON(http.StatusOK, user.Todos)
}
