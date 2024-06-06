//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type RoomRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.Room, error)
	Get(ctx context.Context, id string) (*entity.Room, error)
	Create(ctx context.Context, room entity.Room) error
	BatchCreate(ctx context.Context, rooms []entity.Room) error
	Update(ctx context.Context, id string, room entity.Room) error
	Delete(ctx context.Context, id string) error
	CreateOrUpdate(ctx context.Context, id string, qcs []QueryCondition, room entity.Room) error
}
