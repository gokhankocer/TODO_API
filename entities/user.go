package entities

import (
	//"github.com/gokhankocer/TODO-API/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	//"time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"name"`
	Email    string `gorm:"email"`
	Password string `gorm:"password"`
	Todos    []Todo `gorm:"foreignKey:UserID"`
	IsActive bool   `gorm:"is_active"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) VerifyPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

/*func (user *User) BeforeSave(*gorm.DB) error {
	if user.Password != "" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			return err
		}
		user.Password = string(bytes)
	}
	return nil
}*/
