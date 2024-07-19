package entity

import (
	"fmt"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type MembershipChannel struct {
	MembershipID string `json:"membership_id" db:"membership_id"`
	ChannelID    string `json:"channel_id" db:"channel_id"`
}

func NewMembershipChannel(membershipID, channelID string) (*MembershipChannel, error) {
	if membershipID == "" {
		log.Warn("MembershipID is required", log.Fstring("membershipID", membershipID))
		return nil, fmt.Errorf("membershipID is required")
	}
	if channelID == "" {
		log.Warn("ChannelID is required", log.Fstring("channelID", channelID))
		return nil, fmt.Errorf("channelID is required")
	}
	return &MembershipChannel{
		MembershipID: membershipID,
		ChannelID:    channelID,
	}, nil
}
