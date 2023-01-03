package entities

import (
	"time"

	"github.com/gokhankocer/TODO-API/database"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Status      string         `json:"status"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created at"`
	UpdatedAt   time.Time      `json:"updated at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted at"`
	UserID      uint
}

func (todo *Todo) Save() (*Todo, error) {
	err := database.DB.Create(&todo).Error
	if err != nil {
		return &Todo{}, err
	}
	return todo, nil
}
