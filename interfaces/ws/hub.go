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
	channels         map[*Channel]bool
	Register         chan *Client
	unregister       chan *Client
	broadcast        chan []byte
	channelUseCase   usecase.ChannelUseCase
	pubsubRepo       repository.PubSubRepository
	messageCacheRepo repository.MessageCacheRepository
}

// NewWebsocketServer creates a new Hub type
func NewHub(name string, channelUseCase usecase.ChannelUseCase, pubsubRepo repository.PubSubRepository, messageCacheRepo repository.MessageCacheRepository) *Hub {
	return &Hub{
		ID:               uuid.New().String(),
		Name:             name,
		clients:          make(map[*Client]bool),
		channels:         make(map[*Channel]bool),
		Register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan []byte),
		channelUseCase:   channelUseCase,
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
	ctx := context.Background()

	h.clients[client] = true

	membershipID := client.UserID + "_" + h.ID
	channels, err := h.channelUseCase.ListMembershipChannels(ctx, membershipID)
	if err != nil {
		log.Error("Failed to list membership channels", log.Fstring("membershipID", membershipID))
		return
	}

	for _, ch := range channels {
		channel := h.FindChannelByID(ch.ID)
		if channel == nil {
			log.Warn("Channel not found", log.Fstring("channelID", ch.ID))
			continue
		}
		client.channels[channel] = true
		channel.register <- client
	}
}

func (h *Hub) unregisterClient(client *Client) {
	for channel := range client.channels {
		channel.unregister <- client
	}
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

func (h *Hub) FindChannelByID(id string) *Channel {
	for channel := range h.channels {
		if channel.ID == id {
			return channel
		}
	}
	return nil
}

func (h *Hub) FindChannelByName(name string) *Channel {
	for channel := range h.channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}

func (h *Hub) CreateChannel(ctx context.Context, membershipID, channelName string, channelPrivate bool) *Channel {
	channel := NewChannel(channelName, channelPrivate, h.pubsubRepo, h.messageCacheRepo)

	if err := h.channelUseCase.CreateChannel(ctx, usecase.CreateChannelParams{
		ID:           channel.ID,
		MembershipID: membershipID,
		WorkspaceID:  h.ID,
		Name:         channel.Name,
		Private:      channel.Private,
	}); err != nil {
		log.Error("Failed to create channel", log.Fstring("name", channelName))
		return nil
	}

	go channel.Run(ctx)
	h.channels[channel] = true
	return channel
}
