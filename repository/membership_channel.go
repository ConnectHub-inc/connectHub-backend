//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type MembershipChannelRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.MembershipChannel, error)
	Get(ctx context.Context, membershipID, channelID string) (*entity.MembershipChannel, error)
	Create(ctx context.Context, membershipChannel entity.MembershipChannel) error
	BatchCreate(ctx context.Context, membershipChannels []entity.MembershipChannel) error
	Update(ctx context.Context, id string, membershipChannel entity.MembershipChannel) error
	Delete(ctx context.Context, membershipID, channelID string) error
}
