package entities

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Status      string         `json:"status"`
	Description string         `json:"description"`
	
}

type User struct {
}
