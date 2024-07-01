package entity

type UserRoom struct {
	UserWorkspaceID string `json:"user_workspace_id" db:"user_workspace_id"`
	RoomID          string `json:"room_id" db:"room_id"`
}

func NewUserRoom(userID, workspaceID, roomID string) UserRoom {
	return UserRoom{
		UserWorkspaceID: userID + "_" + workspaceID,
		RoomID:          roomID,
	}
}
