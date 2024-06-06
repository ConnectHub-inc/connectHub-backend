//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type MessageRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.Message, error)
	Get(ctx context.Context, id string) (*entity.Message, error)
	Create(ctx context.Context, message entity.Message) error
	BatchCreate(ctx context.Context, messages []entity.Message) error
	Update(ctx context.Context, id string, message entity.Message) error
	Delete(ctx context.Context, id string) error
	CreateOrUpdate(ctx context.Context, id string, qcs []QueryCondition, message entity.Message) error
}

type MessageCacheRepository interface {
	Set(ctx context.Context, key string, message entity.Message) error
	Get(ctx context.Context, key string) (*entity.Message, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}
