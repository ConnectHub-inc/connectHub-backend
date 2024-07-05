package entity

import (
	"fmt"
	"strings"
)

type Membership struct {
	UserID          string `json:"user_id" db:"user_id"`
	WorkspaceID     string `json:"workspace_id" db:"workspace_id"`
	Name            string `json:"name" db:"name"`
	ProfileImageURL string `json:"profile_image_url" db:"profile_image_url"`
	IsAdmin         bool   `json:"is_admin" db:"is_admin"`
	IsDeleted       bool   `json:"is_deleted" db:"is_deleted"`
}

func NewMembership(userID, workspaceID, name, profileImageURL string, isAdmin, isDeleted bool) Membership {
	return Membership{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            name,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
		IsDeleted:       isDeleted,
	}
}

func (m *Membership) SplitMembershipID(membershipID string) (string, string, error) {
	parts := strings.Split(membershipID, "_")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid user_workspace_id format")
	}
	return parts[0], parts[1], nil
}
