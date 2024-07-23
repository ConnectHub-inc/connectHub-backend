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
	ID       string
	UserID   string // TODO: membershipIDに変更すべきか？
	conn     *websocket.Conn
	hub      *Hub
	channels map[*Channel]bool
	send     chan []byte
	psr      repository.PubSubRepository
	muc      usecase.MessageUseCase
	mcuc     usecase.MembershipChannelUseCase
}

func NewClient(
	userID string,
	conn *websocket.Conn,
	hub *Hub,
	psr repository.PubSubRepository,
	muc usecase.MessageUseCase,
	mcuc usecase.MembershipChannelUseCase,
) *Client {
	return &Client{
		ID:       uuid.New().String(),
		UserID:   userID,
		conn:     conn,
		hub:      hub,
		channels: make(map[*Channel]bool),
		send:     make(chan []byte, config.ChannelBufferSize),
		psr:      psr,
		muc:      muc,
		mcuc:     mcuc,
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
	case entity.ListMessagesAction:
		client.handleListMessages(ctx, message)
	case entity.CreateMessageAction:
		client.handleCreateMessage(ctx, message)
	case entity.DeleteMessageAction:
		client.handleDeleteMessage(ctx, message)
	case entity.UpdateMessageAction:
		client.handleUpdateMessage(ctx, message)
	case entity.CreatePublicChannelAction:
		client.handleCreatePublicChannel(ctx, message)
	case entity.JoinPublicChannelAction:
		client.handleJoinPublicChannel(ctx, message)
	case entity.LeavePublicChannelAction:
		client.handleLeavePublicChannel(ctx, message)
	default:
		log.Warn("Unknown message action", log.Fstring("action", message.Action))
	}
}

func (client *Client) handleListMessages(ctx context.Context, message entity.WSMessage) {
	channelID := message.TargetID
	start := time.Unix(0, 0)                       // Unixエポックの開始
	end := time.Unix(1<<63-62135596801, 999999999) //nolint:gomnd // Unixエポックの終了
	msgs, err := client.muc.ListMessages(ctx, channelID, start, end)
	if err != nil {
		log.Error("Failed to list messages", log.Ferror(err))
		return
	}

	response := entity.WSMessages{
		Action:   entity.ListMessagesAction,
		TargetID: channelID,
		Contents: msgs,
	}
	client.send <- response.Encode()
}

func (client *Client) handleCreateMessage(ctx context.Context, message entity.WSMessage) {
	channelID := message.TargetID
	message.Content.ID = uuid.New().String()
	message.Content.MembershipID = client.UserID + "_" + client.hub.ID

	if err := client.muc.CreateMessage(ctx, channelID, message.Content); err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}

	if channel := client.hub.FindChannelByID(channelID); channel != nil {
		log.Info("Broadcasting message", log.Fstring("channelID", channelID), log.Fstring("messageID", message.Content.ID))
		channel.broadcast <- &message
	} else {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
	}
}

func (client *Client) handleDeleteMessage(ctx context.Context, message entity.WSMessage) {
	channelID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID

	if err := client.muc.DeleteMessage(ctx, message.Content, membershipID, channelID); err != nil {
		log.Error("Failed to delete message", log.Ferror(err))
		return
	}

	if channel := client.hub.FindChannelByID(channelID); channel != nil {
		log.Info("Broadcasting message", log.Fstring("channelID", channelID), log.Fstring("messageID", message.Content.ID))
		channel.broadcast <- &message
	} else {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
	}
}

func (client *Client) handleUpdateMessage(ctx context.Context, message entity.WSMessage) {
	membershipID := client.UserID + "_" + client.hub.ID
	if err := client.muc.UpdateMessage(ctx, message.Content, membershipID); err != nil {
		log.Error("Failed to update message", log.Ferror(err))
		return
	}

	channelID := message.TargetID
	if channel := client.hub.FindChannelByID(channelID); channel != nil {
		log.Info("Broadcasting message", log.Fstring("channelID", channelID), log.Fstring("messageID", message.Content.ID))
		channel.broadcast <- &message
	} else {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
	}
}

func (client *Client) handleCreatePublicChannel(ctx context.Context, message entity.WSMessage) {
	channelName := message.Content.Text
	membershipID := client.UserID + "_" + client.hub.ID
	channel := client.hub.FindChannelByName(channelName)
	if channel != nil {
		log.Warn("Channel already exists", log.Fstring("channelName", channelName))
		return
	}

	channel = client.hub.CreateChannel(ctx, membershipID, channelName, "", false) // TODO: descriptionを追加する
	if channel == nil {
		log.Error("Failed to create channel", log.Fstring("channelName", channelName))
		return
	}

	for c := range client.hub.clients {
		if c.UserID == client.UserID {
			if !c.isInChannel(channel) {
				c.channels[channel] = true
				channel.register <- c
				log.Info("Client registered to channel", log.Fstring("clientID", c.ID), log.Fstring("channelID", channel.ID))
			}
		}
	}

	time.Sleep(5 * time.Second) //nolint:gomnd // TODO: time.Sleepを使うのは避ける

	content, err := entity.NewMessage(membershipID, channel.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(entity.CreatePublicChannelAction, *content, channel.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	channel.broadcast <- msg
}

func (client *Client) handleJoinPublicChannel(ctx context.Context, message entity.WSMessage) {
	channelID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID
	channel := client.hub.FindChannelByID(channelID)
	if channel == nil {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
		return
	}

	if err := client.mcuc.CreateMembershipChannel(ctx, membershipID, channelID); err != nil {
		log.Error("Failed to create membership channel", log.Fstring("membershipID", membershipID), log.Fstring("channelID", channelID))
		return
	}

	for c := range client.hub.clients {
		if c.UserID == client.UserID {
			if !c.isInChannel(channel) {
				c.channels[channel] = true
				channel.register <- c
				log.Info("Client registered to channel", log.Fstring("clientID", c.ID), log.Fstring("channelID", channel.ID))
			}
		}
	}

	content, err := entity.NewMessage(membershipID, channel.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(entity.JoinPublicChannelAction, *content, channel.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	channel.broadcast <- msg
}

func (client *Client) handleLeavePublicChannel(ctx context.Context, message entity.WSMessage) {
	channelID := message.TargetID
	membershipID := client.UserID + "_" + client.hub.ID
	channel := client.hub.FindChannelByID(channelID)
	if channel == nil {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
		return
	}

	if err := client.mcuc.DeleteMembershipChannel(ctx, membershipID, channelID); err != nil {
		log.Error("Failed to delete membership channel", log.Fstring("membershipID", membershipID), log.Fstring("channelID", channelID))
		return
	}

	for c := range client.hub.clients {
		if c.UserID == client.UserID {
			if !c.isInChannel(channel) {
				delete(c.channels, channel)
				channel.unregister <- c
				log.Info("Client unregistered from channel", log.Fstring("clientID", client.ID), log.Fstring("channelID", channel.ID))
			}
		}
	}

	content, err := entity.NewMessage(membershipID, channel.Name)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	msg, err := entity.NewWSMessage(entity.LeavePublicChannelAction, *content, channel.ID, client.ID)
	if err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return
	}
	channel.broadcast <- msg
}

func (client *Client) isInChannel(channel *Channel) bool {
	if _, ok := client.channels[channel]; ok {
		return true
	}
	return false
}
