package ws

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type Channel struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Private      bool   `json:"private"`
	clients      map[*Client]bool
	register     chan *Client
	unregister   chan *Client
	broadcast    chan *entity.WSMessage
	pubsubRepo   repository.PubSubRepository
	msgCacheRepo repository.MessageCacheRepository
}

func NewChannel(
	name string,
	private bool,
	pubsubRepo repository.PubSubRepository,
	msgCacheRepo repository.MessageCacheRepository,
) *Channel {
	return &Channel{
		ID:           uuid.New().String(),
		Name:         name,
		Private:      private,
		clients:      make(map[*Client]bool),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan *entity.WSMessage),
		pubsubRepo:   pubsubRepo,
		msgCacheRepo: msgCacheRepo,
	}
}

func (channel *Channel) Run(ctx context.Context) {
	go channel.subscribeToChannelMessages(ctx)

	for {
		select {
		case client := <-channel.register:
			channel.registerClientInChannel(client)

		case client := <-channel.unregister:
			channel.unregisterClientInChannel(client)

		case message := <-channel.broadcast:
			channel.publishChannelMessage(ctx, message)
		}
	}
}

func (channel *Channel) registerClientInChannel(client *Client) {
	if !channel.Private {
		channel.notifyClientJoined(client)
	}
	channel.clients[client] = true
}

func (channel *Channel) unregisterClientInChannel(client *Client) {
	delete(channel.clients, client)
}

func (channel *Channel) notifyClientJoined(client *Client) {
	membershipID := client.UserID + "_" + client.hub.ID
	content, err := entity.NewMessage(
		membershipID,
		fmt.Sprintf(config.WelcomeMessage, "client.Name"),
	)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}

	message, err := entity.NewWSMessage(
		entity.CreateMessageAction,
		*content,
		channel.ID,
		client.ID,
	)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}

	channel.broadcastToClientsInChannel(message.Encode())
}

func (channel *Channel) broadcastToClientsInChannel(message []byte) {
	for client := range channel.clients {
		client.send <- message
	}
}

func (channel *Channel) publishChannelMessage(ctx context.Context, message *entity.WSMessage) {
	if err := channel.pubsubRepo.Publish(ctx, channel.ID, message.Encode()); err != nil {
		log.Error("Failed to publish message", log.Ferror(err))
	}
}

func (channel *Channel) subscribeToChannelMessages(ctx context.Context) {
	pubsub := channel.pubsubRepo.Subscribe(ctx, channel.ID)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		channel.broadcastToClientsInChannel([]byte(msg.Payload))
	}
}
