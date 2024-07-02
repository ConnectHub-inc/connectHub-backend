//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type UserWorkspaceRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.UserWorkspace, error)
	Get(ctx context.Context, userID, workspaceID string) (*entity.UserWorkspace, error)
	Create(ctx context.Context, userWorkspace entity.UserWorkspace) error
	Update(ctx context.Context, userWorkspace entity.UserWorkspace) error
	Delete(ctx context.Context, ID string) error
}
