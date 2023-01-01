package models

type UserResponse struct {
	ID uint64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Email string `json:"email,omitempty"`
}
