//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/entity"
)

type UserRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.User, error)
	ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.User, error)
	ListRoomUsers(ctx context.Context, channelID string) ([]entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, id string, user entity.User) error
	Delete(ctx context.Context, id string) error
}

type UserCacheRepository interface {
	Set(ctx context.Context, key string, user entity.User) error
	SetUserSession(ctx context.Context, userID string, sessionData string) error
	Get(ctx context.Context, key string) (*entity.User, error)
	GetUserSession(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Scan(ctx context.Context, match string) ([]string, error)
}

type UserRoomRepository interface {
	List(ctx context.Context, qcs []QueryCondition) ([]entity.UserRoom, error)
	Get(ctx context.Context, id string) (*entity.UserRoom, error)
	Create(ctx context.Context, userRoom entity.UserRoom) error
	Update(ctx context.Context, id string, userRoom entity.UserRoom) error
	Delete(ctx context.Context, id string) error
}
