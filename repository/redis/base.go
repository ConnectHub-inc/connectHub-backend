package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

var ErrCacheMiss = errors.New("cache: key not found")

type base[T any] struct {
	client *redis.Client
}

func newBase[T any](client *redis.Client) *base[T] {
	return &base[T]{
		client: client,
	}
}

func (b *base[T]) Set(ctx context.Context, key string, entity T) error {
	serializeEntity, err := b.serialize(entity)
	if err != nil {
		log.Error("Failed to serialize entity", log.Ferror(err))
		return err
	}
	if err = b.client.Set(ctx, key, serializeEntity, 0).Err(); err != nil {
		log.Error("Failed to set cache", log.Ferror(err))
		return err
	}
	log.Info("Cache set successfully", log.Fstring("key", key))
	return nil
}

func (b *base[T]) Get(ctx context.Context, key string) (*T, error) {
	val, err := b.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		log.Warn("Cache miss", log.Fstring("key", key))
		return nil, ErrCacheMiss
	} else if err != nil {
		log.Error("Failed to get cache", log.Ferror(err))
		return nil, err
	}
	entity, err := b.deserialize(val)
	if err != nil {
		log.Error("Failed to deserialize entity", log.Ferror(err))
		return nil, err
	}
	log.Info("Cache hit", log.Fstring("key", key))
	return entity, nil
}

func (b *base[T]) Delete(ctx context.Context, key string) error {
	if err := b.client.Del(ctx, key).Err(); err != nil {
		log.Error("Failed to delete cache", log.Ferror(err))
		return err
	}
	log.Info("Cache deleted successfully", log.Fstring("key", key))
	return nil
}

func (b *base[T]) Exists(ctx context.Context, key string) bool {
	val := b.client.Exists(ctx, key).Val()
	exists := val > 0
	if exists {
		log.Info("Cache exists", log.Fstring("key", key))
	} else {
		log.Warn("Cache does not exist", log.Fstring("key", key))
	}
	return exists
}

func (b *base[T]) Scan(ctx context.Context, match string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, newCursor, err := b.client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			log.Error("Failed to scan cache", log.Ferror(err))
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	log.Info("Cache scan completed", log.Fstring("match", match), log.Fint("keys found", len(allKeys)))
	return allKeys, nil
}

func (b *base[T]) serialize(entity T) (string, error) {
	data, err := json.Marshal(entity)
	if err != nil {
		log.Error("Failed to serialize entity", log.Ferror(err))
		return "", err
	}
	return string(data), nil
}

func (b *base[T]) deserialize(data string) (*T, error) {
	var entity T
	err := json.Unmarshal([]byte(data), &entity)
	if err != nil {
		log.Error("Failed to deserialize entity", log.Ferror(err))
		return nil, err
	}
	return &entity, nil
}
