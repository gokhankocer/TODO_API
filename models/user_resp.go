package models

import "fmt"

type UserResponse struct {
	ID uint64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Email string `json:"email,omitempty"`
}

func (u *UserResponse) Validate() error {
	if u.ID == 0 {
		return fmt.Errorf("invalid ID")
	}
	if u.Name == "" {
		return fmt.Errorf("invalid Name")
	}
	if u.Email == "" {
		return fmt.Errorf("invalid Email")
	}
	return nil
}
