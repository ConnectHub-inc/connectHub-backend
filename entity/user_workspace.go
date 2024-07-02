package entity

type UserWorkspace struct {
	UserID          string `json:"user_id" db:"user_id"`
	WorkspaceID     string `json:"workspace_id" db:"workspace_id"`
	Name            string `json:"name" db:"name"`
	ProfileImageURL string `json:"profile_image_url" db:"profile_image_url"`
	IsAdmin         bool   `json:"is_admin" db:"is_admin"`
}

func NewUserWorkspace(userID, workspaceID, name, profileImageURL string, isAdmin bool) UserWorkspace {
	return UserWorkspace{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            name,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
	}
}
