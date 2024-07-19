package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

const (
	ListMessagesAction        = "LIST_MESSAGES"
	CreateMessageAction       = "CREATE_MESSAGE"
	DeleteMessageAction       = "DELETE_MESSAGE"
	UpdateMessageAction       = "UPDATE_MESSAGE"
	CreatePublicChannelAction = "CREATE_PUBLIC_CHANNEL"
	JoinPublicChannelAction   = "JOIN_PUBLIC_CHANNEL"
	LeavePublicChannelAction  = "LEAVE_PUBLIC_CHANNEL"
)

var validActions = map[string]bool{
	ListMessagesAction:        true,
	CreateMessageAction:       true,
	DeleteMessageAction:       true,
	UpdateMessageAction:       true,
	CreatePublicChannelAction: true,
	JoinPublicChannelAction:   true,
	LeavePublicChannelAction:  true,
}

type Message struct {
	ID           string     `json:"id" db:"id"`
	MembershipID string     `json:"membership_id" db:"membership_id"`
	Text         string     `json:"text" db:"text"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}

type WSMessage struct {
	Action   string  `json:"action_tag"`
	Content  Message `json:"content"`
	TargetID string  `json:"target_id"` // TargetID is the ID of the channel or user the message is intended for
	SenderID string  `json:"sender_id"` // SenderID is the ID of the user who sent the message
}

func (message *WSMessage) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
	}
	return json
}

type WSMessages struct {
	Action   string    `json:"action_tag"`
	Contents []Message `json:"contents"`
	TargetID string    `json:"target_id"` // TargetID is the ID of the channel or user the message is intended for
	SenderID string    `json:"sender_id"` // SenderID is the ID of the user who sent the message
}

func (messages *WSMessages) Encode() []byte {
	json, err := json.Marshal(messages)
	if err != nil {
		log.Error("Failed to encode messages", log.Ferror(err))
	}
	return json
}

func NewMessage(membershipID, text string) (*Message, error) {
	if membershipID == "" || text == "" {
		log.Warn(
			"MembershipID and Text are required",
			log.Fstring("membershipID", membershipID),
			log.Fstring("text", text),
		)
		return nil, fmt.Errorf("id, membershipID and text are required")
	}
	return &Message{
		ID:           uuid.New().String(),
		MembershipID: membershipID,
		Text:         text,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}, nil
}

func NewWSMessage(action string, content Message, targetID, senderID string) (*WSMessage, error) {
	if !validActions[action] {
		log.Warn(
			"Invalid action",
			log.Fstring("action", action),
		)
		return nil, fmt.Errorf("invalid action: %s", action)
	}
	if targetID == "" || senderID == "" {
		log.Warn(
			"TargetID and SenderID are required",
			log.Fstring("targetID", targetID),
			log.Fstring("senderID", senderID),
		)
		return nil, fmt.Errorf("action, targetID and senderID are required")
	}
	return &WSMessage{
		Action:   action,
		Content:  content,
		TargetID: targetID,
		SenderID: senderID,
	}, nil
}

func NewWSMessages(action string, contents []Message, targetID, senderID string) (*WSMessages, error) {
	if !validActions[action] {
		log.Warn(
			"Invalid action",
			log.Fstring("action", action),
		)
		return nil, fmt.Errorf("invalid action: %s", action)
	}
	if len(contents) == 0 || targetID == "" || senderID == "" {
		log.Warn(
			"Contents, TargetID and SenderID are required",
			log.Fstring("targetID", targetID),
			log.Fstring("senderID", senderID),
		)
		return nil, fmt.Errorf("action, contents, targetID and senderID are required")
	}
	return &WSMessages{
		Action:   action,
		Contents: contents,
		TargetID: targetID,
		SenderID: senderID,
	}, nil
}
