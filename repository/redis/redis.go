package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/internal/log"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	conf, err := config.NewCacheConfig(ctx)
	if err != nil || conf == nil {
		log.Error("Failed to load cache config: %s\n", log.Ferror(err))
		return nil
	}

	client := redis.NewClient(&redis.Options{Addr: conf.Addr, Password: conf.Password, DB: conf.DB})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Critical("Failed to connect to Redis", log.Ferror(err), log.Fstring("addr", conf.Addr))
		return nil
	}

	log.Info("Successfully connected to Redis", log.Fstring("addr", conf.Addr))
	return client
}
