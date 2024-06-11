package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MessageUseCase interface {
	CreateMessage(ctx context.Context, message entity.Message) error
	DeleteMessage(ctx context.Context, contentStr string, userID string) error
}

type messageUseCase struct {
	ur  repository.UserRepository
	mr  repository.MessageRepository
	mcr repository.MessageCacheRepository
}

func NewCommentUseCase(
	ur repository.UserRepository,
	mr repository.MessageRepository,
	mcr repository.MessageCacheRepository,
) MessageUseCase {
	return &messageUseCase{
		ur:  ur,
		mr:  mr,
		mcr: mcr,
	}
}

func (muc *messageUseCase) CreateMessage(ctx context.Context, message entity.Message) error {
	if err := muc.mcr.Set(ctx, message.ID, message); err != nil {
		log.Error("Failed to cache message", log.Ferror(err))
		return err
	}
	return nil
}

type DeleteMessageContent struct {
	UserID    string `json:"userID"`
	MessageID string `json:"messageID"`
}

func (muc *messageUseCase) DeleteMessage(ctx context.Context, contentStr string, userID string) error {
	var content DeleteMessageContent

	err := json.Unmarshal([]byte(contentStr), &content)
	if err != nil {
		log.Error("Failed to unmarshal content", log.Fstring("content", contentStr))
		return err
	}

	user, err := muc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Failed to get user", log.Fstring("userID", userID))
		return err
	}

	if !user.IsAdmin && user.ID != content.UserID {
		log.Warn(
			"User don't have permission to delete msg",
			log.Fstring("userID", content.UserID),
			log.Fstring("msgID", content.MessageID),
		)
		return fmt.Errorf("don't have permission to delete msg")
	}

	if err = muc.mcr.Delete(ctx, content.MessageID); err != nil {
		log.Error("Failed to delete msg from cache", log.Fstring("msgID", content.MessageID))
		return err
	}
	return nil
}
