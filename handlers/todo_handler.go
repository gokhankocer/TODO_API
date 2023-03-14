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

type TodoHandler struct {
	TodoRepository repository.TodoRepositoryInterface
	UserRepository repository.UserRepositoryInterface
}

func NewTodoHandler(todoRepository repository.TodoRepositoryInterface, userRepository repository.UserRepositoryInterface) *TodoHandler {
	return &TodoHandler{
		TodoRepository: todoRepository,
		UserRepository: userRepository,
	}
}
func (handler *TodoHandler) AddTodo(c *gin.Context) {
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

	currentUserID, err := helper.CurrentUser(c)
	user, err := handler.UserRepository.GetUserByID(currentUserID)

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
	savedTodo, err := handler.TodoRepository.AddTodo(&todo)

	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": savedTodo})
}

func (handler *TodoHandler) DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := handler.TodoRepository.GetTodo(uint(id))
	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Todo Not Found"})
		return
	}
	currentUserID, err := helper.CurrentUser(c)
	user, err := handler.UserRepository.GetUserByID(currentUserID)

	if err != nil {
		log.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	handler.TodoRepository.DeleteTodo(todo)
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      todo.Status,
		Description: todo.Description,
	}
	c.JSON(http.StatusNoContent, &responseTodo)
}

func (handler *TodoHandler) UpdateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := handler.TodoRepository.GetTodo(uint(id))
	if err != nil {
		log.Println("error", err)

		c.JSON(http.StatusNotFound, gin.H{"message": "Todo Not Found"})
		return
	}
	currentUserID, err := helper.CurrentUser(c)
	user, err := handler.UserRepository.GetUserByID(currentUserID)

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
	handler.TodoRepository.UpdateTodo(todo)
	responseTodo := &models.TodoResponse{
		ID:          uint64(todo.ID),
		Status:      todo.Status,
		Description: todo.Description,
	}
	c.JSON(200, responseTodo)
}

func (handler *TodoHandler) GetTodo(c *gin.Context) {

	currentUserID, err := helper.CurrentUser(c)
	user, err := handler.UserRepository.GetUserByID(currentUserID)

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
	todo, err := handler.TodoRepository.GetTodo(uint(id))
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

func (handler *TodoHandler) GetTodos(c *gin.Context) {
	currentUserID, err := helper.CurrentUser(c)
	user, err := handler.UserRepository.GetUserByID(currentUserID)

	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Error"})
		return
	}

	c.JSON(http.StatusOK, user.Todos)
}
