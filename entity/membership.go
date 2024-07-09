package entity

import (
	"fmt"
	"strings"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Membership struct {
	ID              string `json:"membership_id" db:"id"`
	UserID          string `json:"user_id" db:"user_id"`
	WorkspaceID     string `json:"workspace_id" db:"workspace_id"`
	Name            string `json:"name" db:"name"`
	ProfileImageURL string `json:"profile_image_url" db:"profile_image_url"`
	IsAdmin         bool   `json:"is_admin" db:"is_admin"`
	IsDeleted       bool   `json:"is_deleted" db:"is_deleted"`
}

func NewMembership(userID, workspaceID, name, profileImageURL string, isAdmin bool) (*Membership, error) {
	if userID == "" {
		log.Warn("UserID is required", log.Fstring("userID", userID))
		return nil, fmt.Errorf("userID is required")
	}
	if workspaceID == "" {
		log.Warn("WorkspaceID is required", log.Fstring("workspaceID", workspaceID))
		return nil, fmt.Errorf("workspaceID is required")
	}
	if name == "" {
		log.Warn("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	if profileImageURL == "" {
		profileImageURL = "https://www.hoge.com/avatar.jpg"
	}
	return &Membership{
		ID:              userID + "_" + workspaceID,
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            name,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
		IsDeleted:       false,
	}, nil
}

func (m *Membership) SplitMembershipID(membershipID string) (string, string, error) {
	const expectedParts = 2
	parts := strings.Split(membershipID, "_")
	if len(parts) != expectedParts {
		return "", "", fmt.Errorf("invalid user_workspace_id format")
	}
	return parts[0], parts[1], nil
}
