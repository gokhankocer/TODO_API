package models

type Todo struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
