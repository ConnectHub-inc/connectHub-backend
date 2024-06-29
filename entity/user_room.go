package entity

type UserRoom struct {
	UserID string `json:"user_id" db:"user_id"`
	RoomID string `json:"room_id" db:"room_id"`
}

func NewUserRoom(userID, roomID string) UserRoom {
	return UserRoom{
		UserID: userID,
		RoomID: roomID,
	}
}
