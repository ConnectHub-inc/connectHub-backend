//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type UserWorkspaceUseCase interface {
	ListUsers(ctx context.Context, workspaceID string) ([]entity.UserWorkspace, error)
	GetUser(ctx context.Context, userID, workspaceID string) (*entity.UserWorkspace, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams, user entity.UserWorkspace) error
}

type userWorkspaceUseCase struct {
	uwr repository.UserWorkspaceRepository
}

func NewUserWorkspaceUseCase(uwr repository.UserWorkspaceRepository) UserWorkspaceUseCase {
	return &userWorkspaceUseCase{
		uwr: uwr,
	}
}

func (uwuc *userWorkspaceUseCase) ListUsers(ctx context.Context, workspaceID string) ([]entity.UserWorkspace, error) {
	users, err := uwuc.uwr.List(ctx, []repository.QueryCondition{{Field: "WorkspaceID", Value: workspaceID}})
	if err != nil {
		log.Error("Failed to list workspace users", log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return users, nil
}

func (uwuc *userWorkspaceUseCase) GetUser(ctx context.Context, userID, workspaceID string) (*entity.UserWorkspace, error) {
	user, err := uwuc.uwr.Get(ctx, userID, workspaceID)
	if err != nil {
		log.Error("Failed to get user", log.Fstring("userID", userID), log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return user, nil
}

func (uuc *userWorkspaceUseCase) UpdateUser(ctx context.Context, params *UpdateUserParams, user entity.UserWorkspace) error {
	if user.UserID != params.ID {
		log.Warn(
			"User don't have permission to update user",
			log.Fstring("userID", user.UserID),
			log.Fstring("updateUserID", params.ID),
		)
		return fmt.Errorf("don't have permission to update user")
	}

	user.Name = params.Name
	// user.Email = params.Email
	user.ProfileImageURL = params.ProfileImageURL

	if err := uuc.uwr.Update(ctx, user); err != nil {
		log.Error(
			"Failed to update user",
			log.Fstring("userID", user.UserID),
			log.Fstring("workspaceID", user.WorkspaceID),
		)
		return err
	}
	return nil
}
