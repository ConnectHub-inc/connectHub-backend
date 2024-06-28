package ws

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type Hub struct {
	clients          map[*Client]bool
	rooms            map[*Room]bool
	Register         chan *Client
	unregister       chan *Client
	broadcast        chan []byte
	roomUseCase      usecase.RoomUseCase
	pubsubRepo       repository.PubSubRepository
	messageCacheRepo repository.MessageCacheRepository
}

// NewWebsocketServer creates a new Hub type
func NewHub(pubsubRepo repository.PubSubRepository) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		pubsubRepo: pubsubRepo,
	}
}

// Run starts the server and listens for incoming messages
func (h *Hub) Run() {
	ctx := context.Background()
	go h.listenPubSubChannel(ctx)

	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	delete(h.clients, client)
}

func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.clients {
		client.send <- message
	}
}

func (h *Hub) listenPubSubChannel(ctx context.Context) {
	pubsub := h.pubsubRepo.Subscribe(ctx, config.PubSubGeneralChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		h.broadcastToClients([]byte(msg.Payload))
	}
}

func (h *Hub) FindRoomByID(id string) *Room {
	for room := range h.rooms {
		if room.ID == id {
			return room
		}
	}
	return nil
}

func (h *Hub) FindRoomByName(name string) *Room {
	for room := range h.rooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}

func (h *Hub) CreateRoom(userID, roomName string, roomPrivate bool) *Room {
	room := NewRoom(roomName, roomPrivate, h.pubsubRepo, h.messageCacheRepo)

	if err := h.roomUseCase.CreateRoom(context.Background(), userID, entity.Room{
		ID:      room.ID,
		Name:    room.Name,
		Private: room.Private,
	}); err != nil {
		log.Error("Failed to create room", log.Fstring("name", roomName))
		return nil
	}

	go room.Run(context.Background())
	h.rooms[room] = true
	return room
}
