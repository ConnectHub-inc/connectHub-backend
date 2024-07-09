//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MembershipRoomUseCase interface {
	CreateMembershipRoom(ctx context.Context, membershipID, roomID string) error
	DeleteMembershipRoom(ctx context.Context, membershipID, roomID string) error
}

type membershipRoomUseCase struct {
	mrr repository.MembershipRoomRepository
}

func NewMembershipRoomUseCase(mrr repository.MembershipRoomRepository) MembershipRoomUseCase {
	return &membershipRoomUseCase{
		mrr: mrr,
	}
}

func (mruc *membershipRoomUseCase) CreateMembershipRoom(ctx context.Context, membershipID, roomID string) error {
	membershipRoom, err := entity.NewMembershipRoom(membershipID, roomID)
	if err != nil {
		log.Error("Failed to create membership room", log.Ferror(err))
		return err
	}
	if err = mruc.mrr.Create(ctx, *membershipRoom); err != nil {
		log.Error("Failed to create membership room", log.Fstring("membershipID", membershipID), log.Fstring("roomID", roomID))
		return err
	}
	return nil
}

func (mruc *membershipRoomUseCase) DeleteMembershipRoom(ctx context.Context, membershipID, roomID string) error {
	if err := mruc.mrr.Delete(ctx, membershipID, roomID); err != nil {
		log.Error("Failed to delete membership room", log.Fstring("membershipID", membershipID), log.Fstring("roomID", roomID))
		return err
	}
	return nil
}
