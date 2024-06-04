package entity

import (
	"encoding/json"
	"log"
)

type Message struct {
	Action   string `json:"action"`
	Content  string `json:"content"`
	TargetID string `json:"target"` // TargetID is the ID of the room or user the message is intended for
	SenderID string `json:"sender"` // SenderID is the ID of the user who sent the message
}

func (message *Message) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}
