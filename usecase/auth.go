//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type AuthUseCase interface {
	GetUserFromContext(ctx context.Context) (*entity.User, error)
}

type authUseCase struct {
	ur repository.UserRepository
}

func NewAuthUseCase(ur repository.UserRepository) AuthUseCase {
	return &authUseCase{
		ur: ur,
	}
}

func (auc *authUseCase) GetUserFromContext(ctx context.Context) (*entity.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Printf("Failed to retrieve userId from context")
		return nil, fmt.Errorf("user name not found in request context")
	}
	user, err := auc.ur.Get(ctx, userID)
	if err != nil {
		log.Printf("Failed to get UserInfo from db: %v", userID)
		return nil, err
	}
	return user, nil
}
