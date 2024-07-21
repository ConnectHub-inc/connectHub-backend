//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type ChannelUseCase interface {
	CreateChannel(ctx context.Context, params CreateChannelParams) error
	ListMembershipChannels(ctx context.Context, membershipID string) ([]entity.Channel, error)
}

type channelUseCase struct {
	cr  repository.ChannelRepository
	mrr repository.MembershipChannelRepository
	tr  repository.TransactionRepository
}

func NewChannelUseCase(cr repository.ChannelRepository, mrr repository.MembershipChannelRepository, tr repository.TransactionRepository) ChannelUseCase {
	return &channelUseCase{
		cr:  cr,
		mrr: mrr,
		tr:  tr,
	}
}

type CreateChannelParams struct {
	ID           string
	MembershipID string
	WorkspaceID  string
	Name         string
	Description  string
	Private      bool
}

func (ruc *channelUseCase) CreateChannel(ctx context.Context, params CreateChannelParams) error {
	err := ruc.tr.Transaction(ctx, func(ctx context.Context) error {
		channel, err := entity.NewChannel(params.ID, params.WorkspaceID, params.Name, params.Description, params.Private)
		if err != nil {
			log.Error("Failed to create channel", log.Ferror(err))
			return err
		}
		if err = ruc.cr.Create(ctx, *channel); err != nil {
			log.Error("Failed to create channel", log.Fstring("channelID", channel.ID))
			return err
		}

		var membershipChannel *entity.MembershipChannel
		membershipChannel, err = entity.NewMembershipChannel(params.MembershipID, channel.ID)
		if err != nil {
			log.Error("Failed to create membership channel", log.Ferror(err))
			return err
		}
		if err = ruc.mrr.Create(ctx, *membershipChannel); err != nil {
			log.Error(
				"Failed to create membership channel",
				log.Fstring("membershipID", params.MembershipID),
				log.Fstring("channelID", channel.ID),
			)
			return err
		}

		return nil
	})
	if err != nil {
		log.Error("Failed to create channel", log.Fstring("channelName", params.Name))
		return err
	}
	return nil
}

func (ruc *channelUseCase) ListMembershipChannels(ctx context.Context, membershipID string) ([]entity.Channel, error) {
	channels, err := ruc.cr.ListMembershipChannels(ctx, membershipID)
	if err != nil {
		log.Error("Failed to list user workspace channels", log.Fstring("membershipID", membershipID))
		return nil, err
	}
	return channels, nil
}
