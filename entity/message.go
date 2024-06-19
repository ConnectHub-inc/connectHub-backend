package entity

import (
	"encoding/json"
	"time"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userID"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WSMessage struct {
	Action   string  `json:"action"`
	Content  Message `json:"content"`
	TargetID string  `json:"target"` // TargetID is the ID of the room or user the message is intended for
	SenderID string  `json:"sender"` // SenderID is the ID of the user who sent the message
}

func (message *WSMessage) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
	}
	return json
}

type WSMessages struct {
	Action   string    `json:"action"`
	Contents []Message `json:"content"`
	TargetID string    `json:"target"` // TargetID is the ID of the room or user the message is intended for
	SenderID string    `json:"sender"` // SenderID is the ID of the user who sent the message
}

func (messages *WSMessages) Encode() []byte {
	json, err := json.Marshal(messages)
	if err != nil {
		log.Error("Failed to encode messages", log.Ferror(err))
	}
	return json
}
