//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MembershipChannelUseCase interface {
	CreateMembershipChannel(ctx context.Context, membershipID, channelID string) error
	DeleteMembershipChannel(ctx context.Context, membershipID, channelID string) error
}

type membershipChannelUseCase struct {
	mrr repository.MembershipChannelRepository
}

func NewMembershipChannelUseCase(mrr repository.MembershipChannelRepository) MembershipChannelUseCase {
	return &membershipChannelUseCase{
		mrr: mrr,
	}
}

func (mruc *membershipChannelUseCase) CreateMembershipChannel(ctx context.Context, membershipID, channelID string) error {
	membershipChannel, err := entity.NewMembershipChannel(membershipID, channelID)
	if err != nil {
		log.Error("Failed to create membership channel", log.Ferror(err))
		return err
	}
	if err = mruc.mrr.Create(ctx, *membershipChannel); err != nil {
		log.Error("Failed to create membership channel", log.Fstring("membershipID", membershipID), log.Fstring("channelID", channelID))
		return err
	}
	return nil
}

func (mruc *membershipChannelUseCase) DeleteMembershipChannel(ctx context.Context, membershipID, channelID string) error {
	if err := mruc.mrr.Delete(ctx, membershipID, channelID); err != nil {
		log.Error("Failed to delete membership channel", log.Fstring("membershipID", membershipID), log.Fstring("channelID", channelID))
		return err
	}
	return nil
}
