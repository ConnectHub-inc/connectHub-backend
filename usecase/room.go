//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type RoomUseCase interface {
	CreateRoom(ctx context.Context, params CreateRoomParams) error
	ListMembershipRooms(ctx context.Context, membershipID string) ([]entity.Room, error)
}

type roomUseCase struct {
	rr  repository.RoomRepository
	mrr repository.MembershipRoomRepository
	tr  repository.TransactionRepository
}

func NewRoomUseCase(rr repository.RoomRepository, mrr repository.MembershipRoomRepository, tr repository.TransactionRepository) RoomUseCase {
	return &roomUseCase{
		rr:  rr,
		mrr: mrr,
		tr:  tr,
	}
}

type CreateRoomParams struct {
	ID           string
	MembershipID string
	WorkspaceID  string
	Name         string
	Description  string
	Private      bool
}

func (ruc *roomUseCase) CreateRoom(ctx context.Context, params CreateRoomParams) error {
	err := ruc.tr.Transaction(ctx, func(ctx context.Context) error {
		room, err := entity.NewRoom(params.ID, params.WorkspaceID, params.Name, params.Description, params.Private)
		if err != nil {
			log.Error("Failed to create room", log.Ferror(err))
			return err
		}
		if err = ruc.rr.Create(ctx, *room); err != nil {
			log.Error("Failed to create room", log.Fstring("roomID", room.ID))
			return err
		}

		var membershipRoom *entity.MembershipRoom
		membershipRoom, err = entity.NewMembershipRoom(params.MembershipID, room.ID)
		if err != nil {
			log.Error("Failed to create membership room", log.Ferror(err))
			return err
		}
		if err = ruc.mrr.Create(ctx, *membershipRoom); err != nil {
			log.Error(
				"Failed to create membership room",
				log.Fstring("membershipID", params.MembershipID),
				log.Fstring("roomID", room.ID),
			)
			return err
		}

		return nil
	})
	if err != nil {
		log.Error("Failed to create room", log.Fstring("roomName", params.Name))
		return err
	}
	return nil
}

func (ruc *roomUseCase) ListMembershipRooms(ctx context.Context, membershipID string) ([]entity.Room, error) {
	rooms, err := ruc.rr.ListMembershipRooms(ctx, membershipID)
	if err != nil {
		log.Error("Failed to list user workspace rooms", log.Fstring("membershipID", membershipID))
		return nil, err
	}
	return rooms, nil
}
