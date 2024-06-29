//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type UserRoomRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.UserRoom, error)
	Get(ctx context.Context, userID, roomID string) (*entity.UserRoom, error)
	Create(ctx context.Context, userRoom entity.UserRoom) error
	Update(ctx context.Context, id string, userRoom entity.UserRoom) error
	Delete(ctx context.Context, userID, roomID string) error
}
