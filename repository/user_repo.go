package repository

import (
	"errors"

	"github.com/gokhankocer/TODO-API/database"
	"github.com/gokhankocer/TODO-API/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

type RepositoryInterface interface {
	CreateUser(user *entities.User) error
	HashPassword(password string) (string, error)
}

type Repository struct{}

var Repo RepositoryInterface

func FindUserByName(name string) (entities.User, error) {
	var user entities.User
	err := database.DB.Where("name=?", name).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := database.DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func FindUserById(id uint) (entities.User, error) {
	var user entities.User
	err := database.DB.Preload("Todos").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func GetUserByID(id uint) (*entities.User, error) {
	var user entities.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := database.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *entities.User) error {
	var existingUser entities.User
	err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return errors.New("email already in use")
	}

	return database.DB.Create(user).Error
}

func UpdateUser(id uint, user *entities.User) error {
	return database.DB.Model(&entities.User{}).Where("id = ?", id).Updates(user).Error
}

func DeleteUser(id uint) error {
	return database.DB.Delete(&entities.User{}, id).Error
}
func FindUserByResetPasswordToken(resetPasswordToken string) (*entities.User, error) {
	var user entities.User
	err := database.DB.Where("reset_password_token = ?", resetPasswordToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CheckUserExists(email string) (bool, error) {
	var count int64
	err := database.DB.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *Repository) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
