package models

type UserRequest struct {
	Name string `json:"name,omitempty"`

	Email string `json:"email,omitempty"`
}
