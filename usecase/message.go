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
	UpdateMessage(ctx context.Context, message entity.Message, membershipID string) error
	DeleteMessage(ctx context.Context, message entity.Message, membershipID, channelID string) error
}

type messageUseCase struct {
	ur  repository.MembershipRepository
	mr  repository.MessageRepository
	mcr repository.MessageCacheRepository
}

func NewMessageUseCase(
	ur repository.MembershipRepository,
	mr repository.MessageRepository,
	mcr repository.MessageCacheRepository,
) MessageUseCase {
	return &messageUseCase{
		ur:  ur,
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

func (muc *messageUseCase) UpdateMessage(ctx context.Context, message entity.Message, membershipID string) error {
	membership, err := muc.ur.Get(ctx, membershipID)
	if err != nil {
		log.Error("Failed to get membership", log.Fstring("membershipID", membershipID))
		return err
	}

	if !membership.IsAdmin && membershipID != message.MembershipID {
		log.Warn(
			"Membership don't have permission to update msg",
			log.Fstring("membershipID", membershipID),
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

func (muc *messageUseCase) DeleteMessage(ctx context.Context, message entity.Message, membershipID, channelID string) error {
	membership, err := muc.ur.Get(ctx, membershipID)
	if err != nil {
		log.Error("Failed to get membership", log.Fstring("membershipID", membershipID))
		return err
	}

	if !membership.IsAdmin && membershipID != message.MembershipID {
		log.Warn(
			"Membership don't have permission to delete msg",
			log.Fstring("membershipID", membershipID),
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
