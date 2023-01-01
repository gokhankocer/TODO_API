package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"github.com/gokhankocer/TODO-API/helper"
	"github.com/gokhankocer/TODO-API/models"
	//"github.com/gokhankocer/TODO-API/repository"
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

	/*err := requestTodo.Validate(strfmt.Default)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing Fields",
		})
		return
	}*/
	user, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Logged in"})
		return
	}
	var todo = entities.Todo{
		Status:      requestTodo.Status,
		Description: requestTodo.Description,
	}
	todo.UserID = user.ID
	savedTodo, err := todo.Save()
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create todo"})
		return
	}
	fmt.Println(err)
	c.JSON(http.StatusCreated, gin.H{"data": savedTodo})
}

func DeleteTodo(c *gin.Context) {
	var todo entities.Todo
	if err := database.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Todo Not Found"})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	database.DB.Delete(&todo)
	c.JSON(200, &todo)
}

func UpdateTodo(c *gin.Context) {
	var todo entities.Todo
	if err := database.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Todo Not Found"})
		return
	}
	user, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}
	c.BindJSON(&todo)
	database.DB.Save(&todo)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	var todo entities.Todo
	if err := database.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Todo not found"})
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
		return
	}

	c.JSON(http.StatusAccepted, &todo)

}

func GetTodos(c *gin.Context) {
	user, err := helper.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Error"})
		return
	}
	c.JSON(http.StatusOK, user.Todos)
}