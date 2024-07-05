package entity

import (
	"fmt"
	"strings"
)

type Membership struct {
	ID              string `json:"id" db:"id"`
	UserID          string `json:"user_id" db:"user_id"`
	WorkspaceID     string `json:"workspace_id" db:"workspace_id"`
	Name            string `json:"name" db:"name"`
	ProfileImageURL string `json:"profile_image_url" db:"profile_image_url"`
	IsAdmin         bool   `json:"is_admin" db:"is_admin"`
	IsDeleted       bool   `json:"is_deleted" db:"is_deleted"`
}

func NewMembership(userID, workspaceID, name, profileImageURL string, isAdmin, isDeleted bool) Membership {
	return Membership{
		ID:              userID + "_" + workspaceID,
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            name,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
		IsDeleted:       isDeleted,
	}
}

func (m *Membership) SplitMembershipID(membershipID string) (string, string, error) {
	const expectedParts = 2
	parts := strings.Split(membershipID, "_")
	if len(parts) != expectedParts {
		return "", "", fmt.Errorf("invalid user_workspace_id format")
	}
	return parts[0], parts[1], nil
}
