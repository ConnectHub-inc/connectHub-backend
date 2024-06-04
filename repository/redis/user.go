package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type userRepository struct {
	*base[entity.User]
}

func NewUserRepository(client *redis.Client) repository.UserCacheRepository {
	return &userRepository{
		base: newBase[entity.User](client),
	}
}

func (ur *userRepository) GetUserSession(ctx context.Context, userID string) (string, error) {
	return ur.client.Get(ctx, userID).Result()
}

func (ur *userRepository) SetUserSession(ctx context.Context, userID string, sessionData string) error {
	return ur.client.Set(ctx, userID, sessionData, 0).Err()
}
