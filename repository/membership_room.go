//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type MembershipRoomRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.MembershipRoom, error)
	Get(ctx context.Context, membershipID, roomID string) (*entity.MembershipRoom, error)
	Create(ctx context.Context, membershipRoom entity.MembershipRoom) error
	Update(ctx context.Context, id string, membershipRoom entity.MembershipRoom) error
	Delete(ctx context.Context, membershipID, roomID string) error
}
