package models

import (
	"time"

	"gorm.io/gorm"
)

type Todos struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Status      *string        `json:"status"`
	Description *string        `json:"description"`
	CreatedAt   time.Time      `json:"created at"`
	UpdatedAt   time.Time      `json:"updated at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted at"`
}
