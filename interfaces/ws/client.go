package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

var newline = []byte{'\n'}

type Client struct {
	ID         string
	Name       string
	conn       *websocket.Conn
	hub        *Hub
	rooms      map[*Room]bool
	send       chan []byte
	pubsubRepo repository.PubSubRepository
}

func NewClient(
	conn *websocket.Conn,
	hub *Hub,
	pubsubRepo repository.PubSubRepository,
) *Client {
	return &Client{
		conn:       conn,
		hub:        hub,
		send:       make(chan []byte, config.ChannelBufferSize),
		pubsubRepo: pubsubRepo,
	}
}

func (client *Client) ReadPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(config.MaxMessageSize)
	if err := client.conn.SetReadDeadline(time.Now().Add(config.PongWait)); err != nil {
		log.Printf("failed to set read deadline: %v", err)
	}
	client.conn.SetPongHandler(func(string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(config.PongWait))
		if err != nil {
			log.Printf("Error setting read deadline: %v", err)
			return err
		}
		return nil
	})
	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
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
				return
			}
			if !ok {
				// The Hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{}) //nolint: errcheck // Ignore error.
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err = w.Write(message); err != nil {
				return
			}
			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					return
				}
				if _, err = w.Write(<-client.send); err != nil {
					return
				}
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.conn.SetWriteDeadline(time.Now().Add(config.WriteWait)); err != nil {
				return
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client
	close(client.send)
	client.conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message entity.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshalling JSON message: %v", err)
		return
	}

	message.SenderID = client.ID // TODO: clientのIDでなく、userのIDに変更する

	switch message.Action {
	case config.SendMessageAction:
		roomID := message.TargetID
		if room := client.hub.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}
	case config.CreateRoomAction:
		client.handleCreateRoomMessage(message)
	}
}

func (client *Client) handleCreateRoomMessage(message entity.Message) {
	room := client.hub.createRoom(message.Content, false)
	if room == nil {
		return
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}
