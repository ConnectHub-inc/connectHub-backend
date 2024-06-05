package websocket

import (
	"context"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[*Room]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	pubsubRepo repository.PubSubRepository
	roomRepo   repository.RoomRepository
}

// NewWebsocketServer creates a new Hub type
func NewHub(pubsubRepo repository.PubSubRepository) repository.HubWebSocketRepository {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
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
		case client := <-h.register:
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

func (h *Hub) findRoomByID(id string) *Room {
	for room := range h.rooms {
		if room.ID == id {
			return room
		}
	}
	return nil
}

func (h *Hub) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private, h.pubsubRepo)

	if err := h.roomRepo.Create(context.Background(), entity.Room{
		ID:      room.ID,
		Name:    room.Name,
		Private: room.Private,
	}); err != nil {
		return nil
	}

	go room.Run(context.Background())
	h.rooms[room] = true
	return room
}
