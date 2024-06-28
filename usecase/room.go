//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type RoomUseCase interface {
	CreateRoom(ctx context.Context, userID string, room entity.Room) error
	ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error)
}

type roomUseCase struct {
	rr  repository.RoomRepository
	urr repository.UserRoomRepository
	tr  repository.TransactionRepository
}

func NewRoomUseCase(rr repository.RoomRepository, urr repository.UserRoomRepository, tr repository.TransactionRepository) RoomUseCase {
	return &roomUseCase{
		rr:  rr,
		urr: urr,
		tr:  tr,
	}
}

func (ruc *roomUseCase) CreateRoom(ctx context.Context, userID string, room entity.Room) error {
	err := ruc.tr.Transaction(ctx, func(ctx context.Context) error {
		if err := ruc.rr.Create(ctx, room); err != nil {
			log.Error("Failed to create room", log.Fstring("roomID", room.ID))
			return err
		}

		if err := ruc.urr.Create(ctx, entity.UserRoom{
			UserID: userID,
			RoomID: room.ID,
		}); err != nil {
			log.Error("Failed to create user room", log.Fstring("userID", userID), log.Fstring("roomID", room.ID))
			return err
		}

		return nil
	})
	if err != nil {
		log.Error("Failed to create room", log.Fstring("roomID", room.ID))
		return err
	}
	return nil
}

func (ruc *roomUseCase) ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error) {
	rooms, err := ruc.rr.ListUserWorkspaceRooms(ctx, userID, workspaceID)
	if err != nil {
		log.Error("Failed to list user workspace rooms", log.Fstring("userID", userID), log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return rooms, nil
}
