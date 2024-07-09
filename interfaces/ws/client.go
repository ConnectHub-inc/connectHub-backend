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
	UserID string // TODO: membershipIDに変更すべきか？
	conn   *websocket.Conn
	hub    *Hub
	rooms  map[*Room]bool
	send   chan []byte
	psr    repository.PubSubRepository
	muc    usecase.MessageUseCase
	mruc   usecase.MembershipRoomUseCase
}

func NewClient(
	userID string,
	conn *websocket.Conn,
	hub *Hub,
	psr repository.PubSubRepository,
	muc usecase.MessageUseCase,
	mruc usecase.MembershipRoomUseCase,
) *Client {
	return &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		conn:   conn,
		hub:    hub,
		rooms:  make(map[*Room]bool),
		send:   make(chan []byte, config.ChannelBufferSize),
		psr:    psr,
		muc:    muc,
		mruc:   mruc,
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
	ctx := context.Background()

	var message entity.WSMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Error("Error unmarshalling JSON message", log.Ferror(err))
		return
	}

	message.SenderID = client.ID

	switch message.Action {
	case config.ListMessagesAction:
		client.handleListMessages(ctx, message)
	case config.CreateMessageAction:
		client.handleCreateMessage(ctx, message)
	case config.DeleteMessageAction:
		client.handleDeleteMessage(ctx, message)
	case config.UpdateMessageAction:
		client.handleUpdateMessage(ctx, message)
	case config.CreatePublicRoomAction:
		client.handleCreatePublicRoom(ctx, message)
	case config.JoinPublicRoomAction:
		client.handleJoinPublicRoom(ctx, message)
	case config.LeavePublicRoomAction:
		client.handleLeavePublicRoom(ctx, message)
	default:
		log.Warn("Unknown message action", log.Fstring("action", message.Action))
	}
}

func (client *Client) handleListMessages(ctx context.Context, message entity.WSMessage) {
	roomID := message.TargetID
	start := time.Unix(0, 0)                       // Unixエポックの開始
	end := time.Unix(1<<63-62135596801, 999999999) //nolint:gomnd // Unixエポックの終了
	msgs, err := client.muc.ListMessages(ctx, roomID, start, end)
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

func (client *Client) handleCreateMessage(ctx context.Context, message entity.WSMessage) {
	roomID := message.TargetID
	message.Content.ID = uuid.New().String()
	message.Content.MembershipID = client.UserID + "_" + client.hub.ID

	if err := client.muc.CreateMessage(ctx, roomID, message.Content); err != nil {
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

func (client *Client) handleDeleteMessage(ctx context.Context, message entity.WSMessage) {
	roomID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID

	if err := client.muc.DeleteMessage(ctx, message.Content, membershipID, roomID); err != nil {
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

func (client *Client) handleUpdateMessage(ctx context.Context, message entity.WSMessage) {
	membershipID := client.UserID + "_" + client.hub.ID
	if err := client.muc.UpdateMessage(ctx, message.Content, membershipID); err != nil {
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

func (client *Client) handleCreatePublicRoom(ctx context.Context, message entity.WSMessage) {
	roomName := message.Content.Text
	membershipID := client.UserID + "_" + client.hub.ID
	room := client.hub.FindRoomByName(roomName)
	if room != nil {
		log.Warn("Room already exists", log.Fstring("roomName", roomName))
		return
	}

	room = client.hub.CreateRoom(ctx, membershipID, roomName, false)
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

	time.Sleep(5 * time.Second) //nolint:gomnd // TODO: time.Sleepを使うのは避ける

	content, err := entity.NewMessage(membershipID, room.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(config.CreatePublicRoomAction, *content, room.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	room.broadcast <- msg
}

func (client *Client) handleJoinPublicRoom(ctx context.Context, message entity.WSMessage) {
	roomID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID
	room := client.hub.FindRoomByID(roomID)
	if room == nil {
		log.Warn("Room not found", log.Fstring("roomID", roomID))
		return
	}

	if err := client.mruc.CreateMembershipRoom(ctx, membershipID, roomID); err != nil {
		log.Error("Failed to create membership room", log.Fstring("membershipID", membershipID), log.Fstring("roomID", roomID))
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

	content, err := entity.NewMessage(membershipID, room.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(config.JoinPublicRoomAction, *content, room.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	room.broadcast <- msg
}

func (client *Client) handleLeavePublicRoom(ctx context.Context, message entity.WSMessage) {
	roomID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID
	room := client.hub.FindRoomByID(roomID)
	if room == nil {
		log.Warn("Room not found", log.Fstring("roomID", roomID))
		return
	}

	if err := client.mruc.DeleteMembershipRoom(ctx, membershipID, roomID); err != nil {
		log.Error("Failed to delete membership room", log.Fstring("membershipID", membershipID), log.Fstring("roomID", roomID))
		return
	}

	for c := range client.hub.clients {
		if c.UserID == client.UserID {
			if !c.isInRoom(room) {
				delete(c.rooms, room)
				room.unregister <- c
				log.Info("Client unregistered from room", log.Fstring("clientID", client.ID), log.Fstring("roomID", room.ID))
			}
		}
	}

	content, err := entity.NewMessage(membershipID, room.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(config.LeavePublicRoomAction, *content, room.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	room.broadcast <- msg
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}
