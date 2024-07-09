package redis

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
)

func Test_MessageRepository(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New().String()
	membershipID := uuid.New().String()
	msgs := []entity.Message{
		{ID: uuid.New().String(), MembershipID: membershipID, Text: "content1", CreatedAt: time.Now(), UpdatedAt: nil},
		{ID: uuid.New().String(), MembershipID: membershipID, Text: "content2", CreatedAt: time.Now(), UpdatedAt: nil},
		{ID: uuid.New().String(), MembershipID: membershipID, Text: "content3", CreatedAt: time.Now(), UpdatedAt: nil},
	}

	repo := NewMessageRepository(client)

	// Create messages
	for _, msg := range msgs {
		if err := repo.Create(ctx, channelID, msg); err != nil {
			t.Errorf("Failed to create message: %v", err)
		}
	}

	// Get messages
	getMsg, err := repo.Get(ctx, msgs[0].ID)
	ValidateErr(t, err, nil)
	if getMsg.ID != msgs[0].ID {
		t.Errorf("Expected message ID %s, got %s", msgs[0].ID, getMsg.ID)
	}

	// List messages
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)
	getMsgs, err := repo.List(ctx, channelID, start, end)
	ValidateErr(t, err, nil)
	if len(getMsgs) != len(msgs) {
		t.Errorf("Expected %d messages, got %d", len(msgs), len(getMsgs))
	}

	// Update message
	msgs[0].Text = "updated content"
	if err = repo.Update(ctx, msgs[0]); err != nil {
		t.Errorf("Failed to update message: %v", err)
	}
	getMsg, err = repo.Get(ctx, msgs[0].ID)
	ValidateErr(t, err, nil)
	if getMsg.Text != "updated content" {
		t.Errorf("Expected message text 'updated content', got %s", getMsg.Text)
	}

	// Delete message
	if err = repo.Delete(ctx, channelID, msgs[0].ID); err != nil {
		t.Errorf("Failed to delete message: %v", err)
	}
	_, err = repo.Get(ctx, msgs[0].ID)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
