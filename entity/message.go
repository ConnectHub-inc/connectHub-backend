package entity

import (
	"encoding/json"
	"time"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type MessageContent struct {
	UserID    string    `json:"userID"`
	MessageID string    `json:"messageID"`
	Text      string    `json:"text"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type Message struct {
	ID       string         `json:"id"`
	Action   string         `json:"action"`
	Content  MessageContent `json:"content"`
	TargetID string         `json:"target"` // TargetID is the ID of the room or user the message is intended for
	SenderID string         `json:"sender"` // SenderID is the ID of the user who sent the message
}

func (message *Message) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
	}
	return json
}
