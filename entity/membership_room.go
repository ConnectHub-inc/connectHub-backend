package entity

type MembershipRoom struct {
	MembershipID string `json:"membership_id" db:"membership_id"`
	RoomID       string `json:"room_id" db:"room_id"`
}

func NewMembershipRoom(membershipID, roomID string) MembershipRoom {
	return MembershipRoom{
		MembershipID: membershipID,
		RoomID:       roomID,
	}
}
