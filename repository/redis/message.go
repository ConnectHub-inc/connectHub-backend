package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
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

func (mr *messageRepository) Get(ctx context.Context, id string) (*entity.Message, error) {
	var message entity.Message

	messageBytes, err := mr.client.HGet(ctx, "messages", id).Result()
	if err != nil {
		log.Error("Failed to get message from hash", log.Ferror(err))
		return nil, err
	}

	err = json.Unmarshal([]byte(messageBytes), &message)
	if err != nil {
		log.Error("Failed to unmarshal message", log.Ferror(err))
		return nil, err
	}

	return &message, nil
}

func (mr *messageRepository) List(ctx context.Context, channelID string, start, end time.Time) ([]entity.Message, error) { //nolint:lll // Ignore long line length
	var messages []entity.Message

	messageIDs, err := mr.client.ZRangeByScore(ctx, channelID, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", start.Unix()),
		Max: fmt.Sprintf("%d", end.Unix()),
	}).Result()
	if err != nil {
		log.Error("Failed to get message IDs from sorted set", log.Ferror(err))
		return nil, err
	}

	for _, id := range messageIDs {
		var messageBytes string
		messageBytes, err = mr.client.HGet(ctx, "messages", id).Result()
		if err != nil {
			log.Error("Failed to get message from hash", log.Ferror(err))
			return nil, err
		}

		var message entity.Message
		if err = json.Unmarshal([]byte(messageBytes), &message); err != nil {
			log.Error("Failed to unmarshal message", log.Ferror(err))
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (mr *messageRepository) Create(ctx context.Context, channelID string, message entity.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to serialize message", log.Ferror(err))
		return err
	}

	pipe := mr.client.TxPipeline()

	pipe.HSet(ctx, "messages", message.ID, messageBytes)

	pipe.ZAdd(ctx, channelID, &redis.Z{
		Score:  float64(message.CreatedAt.Unix()),
		Member: message.ID,
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return err
	}
	return nil
}

func (mr *messageRepository) Update(ctx context.Context, message entity.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to serialize message", log.Ferror(err))
		return err
	}

	if err = mr.client.HSet(ctx, "messages", message.ID, messageBytes).Err(); err != nil {
		log.Error("Failed to update message in hash", log.Ferror(err))
		return err
	}

	return nil
}

func (mr *messageRepository) Delete(ctx context.Context, channelID, messageID string) error {
	pipe := mr.client.TxPipeline()

	pipe.HDel(ctx, "messages", messageID)

	pipe.ZRem(ctx, channelID, messageID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error("Failed to delete message", log.Ferror(err))
		return err
	}

	return nil
}
