package repository

import (
	"github.com/gokhankocer/TODO-API/entities"
	"gorm.io/gorm"
)

type todoRepo struct {
	DB *gorm.DB
}
type TodoRepositoryInterface interface {
	AddTodo(t *entities.Todo) (*entities.Todo, error)
	GetTodos(t []*entities.Todo) ([]*entities.Todo, error)
	GetTodo(id uint) (*entities.Todo, error)
	DeleteTodo(t *entities.Todo) error
	UpdateTodo(todo *entities.Todo) error
}

func NewTodoRepository(DB *gorm.DB) TodoRepositoryInterface {
	return &todoRepo{DB}
}

func (tr *todoRepo) GetTodos(t []*entities.Todo) ([]*entities.Todo, error) {
	err := tr.DB.Find(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tr *todoRepo) AddTodo(t *entities.Todo) (*entities.Todo, error) {
	err := tr.DB.Create(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tr *todoRepo) GetTodo(id uint) (*entities.Todo, error) {
	todo := &entities.Todo{ID: id}
	err := tr.DB.First(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (tr *todoRepo) DeleteTodo(t *entities.Todo) error {
	err := tr.DB.Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *todoRepo) UpdateTodo(todo *entities.Todo) error {
	return tr.DB.Save(todo).Error
}

func (tr *todoRepo) GetTodoByID(id uint) (*entities.Todo, error) {
	var todo entities.Todo
	err := tr.DB.Where("id = ?", id).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}
