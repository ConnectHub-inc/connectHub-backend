package entity

import (
	"fmt"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Room struct {
	ID          string `json:"id" db:"id"`
	WorkspaceID string `json:"workspace_id" db:"workspace_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Private     bool   `json:"private" db:"private"`
}

func NewRoom(id, workspaceID, name, description string, private bool) (*Room, error) {
	if id == "" {
		log.Warn("ID is required", log.Fstring("id", id))
		return nil, fmt.Errorf("id is required")
	}
	if workspaceID == "" {
		log.Warn("WorkspaceID is required", log.Fstring("workspaceID", workspaceID))
		return nil, fmt.Errorf("workspaceID is required")
	}
	if name == "" {
		log.Warn("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	return &Room{
		ID:          id,
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		Private:     private,
	}, nil
}
