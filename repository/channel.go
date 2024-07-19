//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type ChannelRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.Channel, error)
	ListMembershipChannels(ctx context.Context, membershipID string) ([]entity.Channel, error)
	Get(ctx context.Context, id string) (*entity.Channel, error)
	Create(ctx context.Context, channel entity.Channel) error
	BatchCreate(ctx context.Context, channels []entity.Channel) error
	Update(ctx context.Context, id string, channel entity.Channel) error
	Delete(ctx context.Context, id string) error
	CreateOrUpdate(ctx context.Context, id string, qcs []QueryCondition, channel entity.Channel) error
}
