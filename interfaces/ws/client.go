package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

var newline = []byte{'\n'}

type Client struct {
	ID     string
	UserID string
	Name   string
	conn   *websocket.Conn
	hub    *Hub
	rooms  map[*Room]bool
	send   chan []byte
	psr    repository.PubSubRepository
	muc    usecase.MessageUseCase
}

func NewClient(
	userID string,
	name string,
	conn *websocket.Conn,
	hub *Hub,
	psr repository.PubSubRepository,
	muc usecase.MessageUseCase,
) *Client {
	return &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   name,
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, config.ChannelBufferSize),
		psr:    psr,
		muc:    muc,
	}
}

func (client *Client) ReadPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(config.MaxMessageSize)
	if err := client.conn.SetReadDeadline(time.Now().Add(config.PongWait)); err != nil {
		log.Error("Failed to set read deadline", log.Ferror(err))
	}
	client.conn.SetPongHandler(func(string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(config.PongWait))
		if err != nil {
			log.Error("Error setting read deadline", log.Ferror(err))
			return err
		}
		return nil
	})
	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn("Unexpected close error", log.Ferror(err))
			} else {
				log.Info("Client disconnected", log.Ferror(err))
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) WritePump() { //nolint: gocognit
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if err := client.conn.SetWriteDeadline(time.Now().Add(config.WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if !ok {
				// The Hub closed the channel.
				if err := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Warn("Failed to write close message", log.Ferror(err))
				}
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("Failed to get next writer", log.Ferror(err))
				return
			}

			if _, err = w.Write(message); err != nil {
				log.Error("Failed to write message", log.Ferror(err))
				return
			}

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					log.Error("Failed to write newline", log.Ferror(err))
					return
				}
				if _, err = w.Write(<-client.send); err != nil {
					log.Error("Failed to write queued message", log.Ferror(err))
					return
				}
			}

			if err = w.Close(); err != nil {
				log.Error("Failed to close writer", log.Ferror(err))
				return
			}
		case <-ticker.C:
			if err := client.conn.SetWriteDeadline(time.Now().Add(config.WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("Failed to write ping message", log.Ferror(err))
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client
	close(client.send)
	if err := client.conn.Close(); err != nil {
		log.Warn("Failed to close connection", log.Ferror(err))
	} else {
		log.Info("Client disconnected successfully", log.Fstring("clientID", client.ID))
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message entity.WSMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Error("Error unmarshalling JSON message", log.Ferror(err))
		return
	}

	message.SenderID = client.ID

	switch message.Action {
	case config.ListMessagesAction:
		client.handleListMessages(message)
	case config.CreateMessageAction:
		client.handleCreateMessage(message)
	case config.DeleteMessageAction:
		client.handleDeleteMessage(message)
	case config.UpdateMessageAction:
		client.handleUpdateMessage(message)
	case config.CreatePublicRoomAction:
		client.handleCreatePublicRoom(message)
	default:
		log.Warn("Unknown message action", log.Fstring("action", message.Action))
	}
}

func (client *Client) handleListMessages(message entity.WSMessage) {
	roomID := message.TargetID
	start := time.Unix(0, 0)                       // Unixエポックの開始
	end := time.Unix(1<<63-62135596801, 999999999) //nolint:gomnd // Unixエポックの終了
	msgs, err := client.muc.ListMessages(context.Background(), roomID, start, end)
	if err != nil {
		log.Error("Failed to list messages", log.Ferror(err))
		return
	}

	response := entity.WSMessages{
		Action:   config.ListMessagesAction,
		TargetID: roomID,
		Contents: msgs,
	}
	client.send <- response.Encode()
}

func (client *Client) handleCreateMessage(message entity.WSMessage) {
	roomID := message.TargetID

	if err := client.muc.CreateMessage(context.Background(), roomID, message.Content); err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}

	if room := client.hub.FindRoomByID(roomID); room != nil {
		log.Info("Broadcasting message", log.Fstring("roomID", roomID), log.Fstring("messageID", message.Content.ID))
		room.broadcast <- &message
	} else {
		log.Warn("Room not found", log.Fstring("roomID", roomID))
	}
}

func (client *Client) handleDeleteMessage(message entity.WSMessage) {
	roomID := message.TargetID

	if err := client.muc.DeleteMessage(context.Background(), message.Content, roomID, client.ID); err != nil {
		log.Error("Failed to delete message", log.Ferror(err))
		return
	}

	if room := client.hub.FindRoomByID(roomID); room != nil {
		log.Info("Broadcasting message", log.Fstring("roomID", roomID), log.Fstring("messageID", message.Content.ID))
		room.broadcast <- &message
	} else {
		log.Warn("Room not found", log.Fstring("roomID", roomID))
	}
}

func (client *Client) handleUpdateMessage(message entity.WSMessage) {
	if err := client.muc.UpdateMessage(context.Background(), message.Content, client.ID); err != nil {
		log.Error("Failed to update message", log.Ferror(err))
		return
	}

	roomID := message.TargetID
	if room := client.hub.FindRoomByID(roomID); room != nil {
		log.Info("Broadcasting message", log.Fstring("roomID", roomID), log.Fstring("messageID", message.Content.ID))
		room.broadcast <- &message
	} else {
		log.Warn("Room not found", log.Fstring("roomID", roomID))
	}
}

func (client *Client) handleCreatePublicRoom(message entity.WSMessage) {
	roomName := message.Content.Text
	room := client.hub.FindRoomByName(roomName)
	if room != nil {
		log.Warn("Room already exists", log.Fstring("roomName", roomName))
		return
	}

	room = client.hub.CreateRoom(roomName, false)
	if room == nil {
		log.Error("Failed to create room", log.Fstring("roomName", roomName))
		return
	}

	for c := range client.hub.clients {
		if c.UserID == client.UserID {
			if !c.isInRoom(room) {
				c.rooms[room] = true
				room.register <- c
				log.Info("Client registered to room", log.Fstring("clientID", c.ID), log.Fstring("roomID", room.ID))
			}
		}
	}

	room.broadcast <- &entity.WSMessage{
		Action:   config.CreatePublicRoomAction,
		TargetID: room.ID,
		SenderID: client.ID,
		Content: entity.Message{
			ID:        uuid.New().String(),
			UserID:    client.UserID,
			Text:      config.WelcomeMessage,
			CreatedAt: time.Now(),
		},
	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}
