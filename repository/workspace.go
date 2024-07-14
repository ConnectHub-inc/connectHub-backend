//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type WorkspaceRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.Workspace, error)
	Get(ctx context.Context, id string) (*entity.Workspace, error)
	Create(ctx context.Context, workspace entity.Workspace) error
	Update(ctx context.Context, id string, workspace entity.Workspace) error
	Delete(ctx context.Context, id string) error
}
