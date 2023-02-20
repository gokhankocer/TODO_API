package models

import (
	"errors"
	"strings"
)

type UserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *UserRequest) Validate() error {
	if u.Name == "" {
		return errors.New("name field is required")
	}
	if u.Email == "" {
		return errors.New("email field is required")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	if u.Password == "" {
		return errors.New("password field is required")
	}

	return nil
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ConfirmResetPasswordRequest struct {
	Password string `json:"password"`
}
