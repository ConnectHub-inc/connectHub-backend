//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type UserRoomUseCase interface {
	CreateUserRoom(ctx context.Context, userID, roomID string) error
}

type userRoomUseCase struct {
	urr repository.UserRoomRepository
}

func NewUserRoomUseCase(urr repository.UserRoomRepository) UserRoomUseCase {
	return &userRoomUseCase{
		urr: urr,
	}
}

func (uruc *userRoomUseCase) CreateUserRoom(ctx context.Context, userID, roomID string) error {
	if err := uruc.urr.Create(ctx, entity.UserRoom{
		UserID: userID,
		RoomID: roomID,
	}); err != nil {
		log.Error("Failed to create user room", log.Fstring("userID", userID), log.Fstring("roomID", roomID))
		return err
	}
	return nil
}
