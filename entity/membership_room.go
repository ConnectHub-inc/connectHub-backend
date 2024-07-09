package entity

import (
	"fmt"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type MembershipRoom struct {
	MembershipID string `json:"membership_id" db:"membership_id"`
	RoomID       string `json:"room_id" db:"room_id"`
}

func NewMembershipRoom(membershipID, roomID string) (*MembershipRoom, error) {
	if membershipID == "" {
		log.Warn("MembershipID is required", log.Fstring("membershipID", membershipID))
		return nil, fmt.Errorf("membershipID is required")
	}
	if roomID == "" {
		log.Warn("RoomID is required", log.Fstring("roomID", roomID))
		return nil, fmt.Errorf("roomID is required")
	}
	return &MembershipRoom{
		MembershipID: membershipID,
		RoomID:       roomID,
	}, nil
}
