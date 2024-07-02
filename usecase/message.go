//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MessageUseCase interface {
	ListMessages(ctx context.Context, channelID string, start, end time.Time) ([]entity.Message, error)
	CreateMessage(ctx context.Context, channelID string, message entity.Message) error
	UpdateMessage(ctx context.Context, message entity.Message, userID, workspaceID string) error
	DeleteMessage(ctx context.Context, message entity.Message, userID, workspaceID, channelID string) error
}

type messageUseCase struct {
	uwr repository.UserWorkspaceRepository
	mr  repository.MessageRepository
	mcr repository.MessageCacheRepository
}

func NewMessageUseCase(
	uwr repository.UserWorkspaceRepository,
	mr repository.MessageRepository,
	mcr repository.MessageCacheRepository,
) MessageUseCase {
	return &messageUseCase{
		uwr: uwr,
		mr:  mr,
		mcr: mcr,
	}
}

func (muc *messageUseCase) ListMessages(ctx context.Context, channelID string, start, end time.Time) ([]entity.Message, error) {
	messages, err := muc.mcr.List(ctx, channelID, start, end)
	if err != nil {
		log.Error("Failed to get messages", log.Ferror(err))
		return nil, err
	}
	return messages, nil
}

func (muc *messageUseCase) CreateMessage(ctx context.Context, channelID string, message entity.Message) error {
	if err := muc.mcr.Create(ctx, channelID, message); err != nil {
		log.Error("Failed to cache message", log.Ferror(err))
		return err
	}
	return nil
}

func (muc *messageUseCase) UpdateMessage(ctx context.Context, message entity.Message, userID, workspaceID string) error {
	user, err := muc.uwr.Get(ctx, userID, workspaceID)
	if err != nil {
		log.Error("Failed to get user", log.Fstring("userID", userID))
		return err
	}

	userWorkspaceID := userID + "_" + workspaceID
	if !user.IsAdmin && userWorkspaceID != message.UserWorkspaceID {
		log.Warn(
			"User don't have permission to update msg",
			log.Fstring("userID", message.UserWorkspaceID),
			log.Fstring("msgID", message.ID),
		)
		return fmt.Errorf("don't have permission to update msg")
	}

	if err = muc.mcr.Update(ctx, message); err != nil {
		log.Error("Failed to update msg in cache", log.Fstring("msgID", message.ID))
		return err
	}
	return nil
}

func (muc *messageUseCase) DeleteMessage(ctx context.Context, message entity.Message, userID, workspaceID, channelID string) error {
	user, err := muc.uwr.Get(ctx, userID, workspaceID)
	if err != nil {
		log.Error("Failed to get user", log.Fstring("userID", userID))
		return err
	}

	userWorkspaceID := userID + "_" + workspaceID
	if !user.IsAdmin && userWorkspaceID != message.UserWorkspaceID {
		log.Warn(
			"User don't have permission to delete msg",
			log.Fstring("userID", message.UserWorkspaceID),
			log.Fstring("msgID", message.ID),
		)
		return fmt.Errorf("don't have permission to delete msg")
	}

	if err = muc.mcr.Delete(ctx, channelID, message.ID); err != nil {
		log.Error("Failed to delete msg from cache", log.Fstring("msgID", message.ID))
		return err
	}
	return nil
}
