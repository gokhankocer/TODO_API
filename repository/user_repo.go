package repository

import (
	"log"

	"github.com/gokhankocer/TODO-API/entities"
	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}
type UserRepositoryInterface interface {
	FindUserByName(name string) (entities.User, error)
	FindUserByEmail(email string) (entities.User, error)
	FindUserById(id uint) (entities.User, error)
	GetUserByID(id uint) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	CreateUser(user *entities.User) error
	UpdateUser(id uint, user *entities.User) error
	DeleteUser(id uint) error
	FindUserByResetPasswordToken(resetPasswordToken string) (*entities.User, error)
}

func NewUserRepository(DB *gorm.DB) UserRepositoryInterface {
	return &userRepo{DB}
}

func (u *userRepo) FindUserByName(name string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("name=?", name).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (u *userRepo) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (u *userRepo) FindUserById(id uint) (entities.User, error) {
	var user entities.User
	err := u.DB.Preload("Todos").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (u *userRepo) GetUserByID(id uint) (*entities.User, error) {
	var user entities.User
	err := u.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := u.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepo) CreateUser(user *entities.User) error {

	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
	}

	return u.DB.Create(user).Error
}

func (u *userRepo) UpdateUser(id uint, user *entities.User) error {
	return u.DB.Model(&entities.User{}).Where("id = ?", id).Updates(user).Error
}

func (u *userRepo) DeleteUser(id uint) error {
	return u.DB.Delete(&entities.User{}, id).Error
}
func (u *userRepo) FindUserByResetPasswordToken(resetPasswordToken string) (*entities.User, error) {
	var user entities.User
	err := u.DB.Where("reset_password_token = ?", resetPasswordToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
