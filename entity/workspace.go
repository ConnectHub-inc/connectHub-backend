package entity

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Workspace struct {
	ID   string `json:"workspace_id" db:"id"`
	Name string `json:"name" db:"name"`
}

func NewWorkspace(name string) (*Workspace, error) {
	if name == "" {
		log.Warn("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	return &Workspace{
		ID:   uuid.New().String(),
		Name: name,
	}, nil
}
