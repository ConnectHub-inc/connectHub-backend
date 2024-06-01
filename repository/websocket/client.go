package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/repository"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn *websocket.Conn
	hub  *Hub
	send chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) repository.ClientWebSocketRepository {
	return &Client{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256),
	}
}

func (client *Client) ReadPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(config.MaxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(config.PongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.hub.broadcast <- jsonMessage
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if !ok {
				// The Hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
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
