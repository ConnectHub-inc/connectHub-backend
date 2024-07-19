//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type MembershipRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.Membership, error)
	ListChannelMemberships(ctx context.Context, channelID string) ([]entity.Membership, error)
	Get(ctx context.Context, id string) (*entity.Membership, error)
	Create(ctx context.Context, membership entity.Membership) error
	Update(ctx context.Context, membership entity.Membership) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}
