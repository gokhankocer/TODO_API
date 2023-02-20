package repository

import (
	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
)

func GetTodos(t []*entities.Todo) ([]*entities.Todo, error) {
	err := database.DB.Find(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func AddTodo(t *entities.Todo) (*entities.Todo, error) {
	err := database.DB.Create(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func GetTodo(id uint) (*entities.Todo, error) {
	todo := &entities.Todo{ID: id}
	err := database.DB.First(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func DeleteTodo(t *entities.Todo) error {
	err := database.DB.Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateTodo(todo *entities.Todo) error {
	return database.DB.Save(todo).Error
}

func GetTodoByID(id uint) (*entities.Todo, error) {
	var todo entities.Todo
	err := database.DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}
