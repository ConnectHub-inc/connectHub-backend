package entity

import "github.com/google/uuid"

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

func NewRoom(name string, description string, private bool) *Room {
	// TODO: validate
	return &Room{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Private:     private,
	}
}
