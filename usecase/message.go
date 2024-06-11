//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MessageUseCase interface {
	CreateMessage(ctx context.Context, message entity.Message) error
	UpdateMessage(ctx context.Context, message entity.Message, userID string) error
	DeleteMessage(ctx context.Context, content entity.MessageContent, userID string) error
}

type messageUseCase struct {
	ur  repository.UserRepository
	mr  repository.MessageRepository
	mcr repository.MessageCacheRepository
}

func NewMessageUseCase(
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

func (muc *messageUseCase) UpdateMessage(ctx context.Context, message entity.Message, userID string) error {
	user, err := muc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Failed to get user", log.Fstring("userID", userID))
		return err
	}

	if !user.IsAdmin && user.ID != message.Content.UserID {
		log.Warn(
			"User don't have permission to update msg",
			log.Fstring("userID", message.Content.UserID),
			log.Fstring("msgID", message.Content.MessageID),
		)
		return fmt.Errorf("don't have permission to update msg")
	}

	if err = muc.mcr.Set(ctx, message.ID, message); err != nil {
		log.Error("Failed to update msg in cache", log.Fstring("msgID", message.ID))
		return err
	}
	return nil
}

func (muc *messageUseCase) DeleteMessage(ctx context.Context, content entity.MessageContent, userID string) error {
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