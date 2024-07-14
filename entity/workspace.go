package entity

import (
	"fmt"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Workspace struct {
	ID   string `json:"workspace_id" db:"id"`
	Name string `json:"name" db:"name"`
}

func NewWorkspace(id, name string) (*Workspace, error) {
	if id == "" {
		log.Warn("ID is required", log.Fstring("id", id))
		return nil, fmt.Errorf("id is required")
	}
	if name == "" {
		log.Warn("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	return &Workspace{
		ID:   id,
		Name: name,
	}, nil
}
