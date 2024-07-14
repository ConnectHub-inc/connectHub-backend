package ws

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type HubManager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) Add(workspaceID string, hub *Hub) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.hubs[workspaceID] = hub
}

func (hm *HubManager) Get(workspaceID string) (*Hub, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	hub, exists := hm.hubs[workspaceID]
	return hub, exists
}

type Hub struct {
	ID               string
	Name             string
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
func NewHub(name string, roomUseCase usecase.RoomUseCase, pubsubRepo repository.PubSubRepository, messageCacheRepo repository.MessageCacheRepository) *Hub {
	return &Hub{
		ID:               uuid.New().String(),
		Name:             name,
		clients:          make(map[*Client]bool),
		rooms:            make(map[*Room]bool),
		Register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan []byte),
		roomUseCase:      roomUseCase,
		pubsubRepo:       pubsubRepo,
		messageCacheRepo: messageCacheRepo,
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
	// TODO: HubにあるPublicなチャンネルにclientを登録する
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	// TODO: Hubにあるチャンネルからclientを削除する
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

func (h *Hub) CreateRoom(ctx context.Context, membershipID, roomName string, roomPrivate bool) *Room {
	room := NewRoom(roomName, roomPrivate, h.pubsubRepo, h.messageCacheRepo)

	if err := h.roomUseCase.CreateRoom(ctx, usecase.CreateRoomParams{
		ID:           room.ID,
		MembershipID: membershipID,
		WorkspaceID:  h.ID,
		Name:         room.Name,
		Private:      room.Private,
	}); err != nil {
		log.Error("Failed to create room", log.Fstring("name", roomName))
		return nil
	}

	go room.Run(ctx)
	h.rooms[room] = true
	return room
}
