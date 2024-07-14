//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type WorkspaceUseCase interface {
	CreateWorkspace(ctx context.Context, name string) error
}

type workspaceUseCase struct {
	wr repository.WorkspaceRepository
}

func NewWorkspaceUseCase(wr repository.WorkspaceRepository) WorkspaceUseCase {
	return &workspaceUseCase{
		wr: wr,
	}
}

func (wuc *workspaceUseCase) CreateWorkspace(ctx context.Context, name string) error {
	workspace, err := entity.NewWorkspace(name)
	if err != nil {
		log.Error("Failed to create workspace", log.Ferror(err))
		return err
	}
	if err = wuc.wr.Create(ctx, *workspace); err != nil {
		log.Error("Failed to create workspace", log.Ferror(err))
		return err
	}
	return nil
}
