package repository

/*import (
	"github.com/gokhankocer/TODO-API/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ToDoRepository struct {
	db *gorm.DB
}

type ToDoRepoInterface interface {
	AddTodo(t *entities.Todo) (*entities.Todo, error)
	DeleteTodo(t *entities.Todo) error
	UpdateTodo(t entities.Todo) error
	GetTodo(id uint) (*entities.Todo, error)
	GetTodos(t []*entities.Todo) ([]*entities.Todo, error)
}

func NewTodoRepository(db *gorm.DB) ToDoRepoInterface {
	return &ToDoRepository{db}
}

func (repository *ToDoRepository) GetTodos(t []*entities.Todo) ([]*entities.Todo, error) {
	err := repository.db.Find(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (repository *ToDoRepository) AddTodo(t *entities.Todo) (*entities.Todo, error) {
	err := repository.db.Create(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (repository *ToDoRepository) GetTodo(id uint) (*entities.Todo, error) {
	todo := &entities.Todo{ID: id}
	err := repository.db.First(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (repository *ToDoRepository) DeleteTodo(t *entities.Todo) error {
	err := repository.db.Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *ToDoRepository) UpdateTodo(t entities.Todo) error {
	return repository.db.Omit(clause.Associations).Save(&t).Error

}*/
