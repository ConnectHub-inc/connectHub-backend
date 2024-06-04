package redis

import (
	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type messageRepository struct {
	*base[entity.Message]
}

func NewMessageRepository(client *redis.Client) repository.MessageCacheRepository {
	return &messageRepository{
		base: newBase[entity.Message](client),
	}
}
