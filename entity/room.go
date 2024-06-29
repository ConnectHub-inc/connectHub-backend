package entity

import "github.com/google/uuid"

type Room struct {
	ID          string `json:"id" db:"id"`
	WorkspaceID string `json:"workspace_id" db:"workspace_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Private     bool   `json:"private" db:"private"`
}

func NewRoom(workspaceID, name, description string, private bool) Room {
	// TODO: validate
	return Room{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		Private:     private,
	}
}
