package entity

import (
	"fmt"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Channel struct {
	ID          string `json:"id" db:"id"`
	WorkspaceID string `json:"workspace_id" db:"workspace_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Private     bool   `json:"private" db:"private"`
}

func NewChannel(id, workspaceID, name, description string, private bool) (*Channel, error) {
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
	return &Channel{
		ID:          id,
		WorkspaceID: workspaceID,
		Name:        name,
		Description: description,
		Private:     private,
	}, nil
}

type DefaultChannel struct {
	Name        string
	Description string
}

var DefaultChannels = []DefaultChannel{
	{
		Name:        "general",
		Description: "This is the general channel where everyone is included.",
	},
	{
		Name:        "random",
		Description: "This is the random channel for non-work-related discussions.",
	},
}
