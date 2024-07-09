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

type Room struct {
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

func NewRoom(
	name string,
	private bool,
	pubsubRepo repository.PubSubRepository,
	msgCacheRepo repository.MessageCacheRepository,
) *Room {
	return &Room{
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

func (room *Room) Run(ctx context.Context) {
	go room.subscribeToRoomMessages(ctx)

	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.publishRoomMessage(ctx, message)
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	if !room.Private {
		room.notifyClientJoined(client)
	}
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	delete(room.clients, client)
}

func (room *Room) notifyClientJoined(client *Client) {
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
		config.CreateMessageAction,
		*content,
		room.ID,
		client.ID,
	)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}

	room.broadcastToClientsInRoom(message.Encode())
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) publishRoomMessage(ctx context.Context, message *entity.WSMessage) {
	if err := room.pubsubRepo.Publish(ctx, room.ID, message.Encode()); err != nil {
		log.Error("Failed to publish message", log.Ferror(err))
	}
}

func (room *Room) subscribeToRoomMessages(ctx context.Context) {
	pubsub := room.pubsubRepo.Subscribe(ctx, room.ID)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}
