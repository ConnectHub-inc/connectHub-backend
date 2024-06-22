//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type RoomUseCase interface {
	ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error)
}

type roomUseCase struct {
	rr repository.RoomRepository
}

func NewRoomUseCase(rr repository.RoomRepository) RoomUseCase {
	return &roomUseCase{
		rr: rr,
	}
}

func (ruc *roomUseCase) ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error) {
	rooms, err := ruc.rr.ListUserWorkspaceRooms(ctx, userID, workspaceID)
	if err != nil {
		log.Error("Failed to list user workspace rooms", log.Fstring("userID", userID), log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return rooms, nil
}
