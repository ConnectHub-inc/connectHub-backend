package entity

import "github.com/google/uuid"

type Room struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func NewRoom(name string, private bool) *Room {
	// TODO: validate
	return &Room{
		ID:      uuid.New().String(),
		Name:    name,
		Private: private,
	}
}
