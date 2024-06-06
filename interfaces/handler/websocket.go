package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/interfaces/ws"
	"github.com/tusmasoma/connectHub-backend/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.BufferSize,
	WriteBufferSize: config.BufferSize,
}

type WebsocketHandler struct {
	pubsubRepo repository.PubSubRepository
}

func NewWebsocketHandler(pubsub repository.PubSubRepository) *WebsocketHandler {
	return &WebsocketHandler{
		pubsubRepo: pubsub,
	}
}

func (wsh *WebsocketHandler) WebSocket(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // conn is *websocket.Conn
	if err != nil {
		log.Println(err)
		return
	}

	client := ws.NewClient(conn, hub, wsh.pubsubRepo)

	go client.WritePump()
	go client.ReadPump()

	hub.Register <- client
}
