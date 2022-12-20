package models

type UserResponse struct {
	ID string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Email string `json:"email,omitempty"`
}
