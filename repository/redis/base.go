package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
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
		return err
	}
	if err = b.client.Set(ctx, key, serializeEntity, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (b *base[T]) Get(ctx context.Context, key string) (*T, error) {
	val, err := b.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	} else if err != nil {
		return nil, err
	}
	entity, err := b.deserialize(val)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (b *base[T]) Delete(ctx context.Context, key string) error {
	err := b.client.Del(ctx, key).Err()
	return err
}

func (b *base[T]) Exists(ctx context.Context, key string) bool {
	val := b.client.Exists(ctx, key).Val()
	return val > 0
}

func (b *base[T]) Scan(ctx context.Context, match string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, newCursor, err := b.client.Scan(ctx, cursor, match, 0).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return allKeys, nil
}

func (b *base[T]) serialize(entity T) (string, error) {
	data, err := json.Marshal(entity)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *base[T]) deserialize(data string) (*T, error) {
	var entity T
	err := json.Unmarshal([]byte(data), &entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
