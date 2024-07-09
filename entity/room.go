package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Room struct {
	ID          string `json:"id" db:"id"`
	WorkspaceID string `json:"workspace_id" db:"workspace_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Private     bool   `json:"private" db:"private"`
}

func NewRoom(workspaceID, name, description string, private bool) (*Room, error) {
	if workspaceID == "" {
		log.Warn("WorkspaceID is required", log.Fstring("workspaceID", workspaceID))
		return nil, fmt.Errorf("workspaceID is required")
	}
	if name == "" {
		log.Warn("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	return &Room{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		Private:     private,
	}, nil
}
