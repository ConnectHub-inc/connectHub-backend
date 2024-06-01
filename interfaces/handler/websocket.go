package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tusmasoma/connectHub-backend/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type WebsocketHandler struct {
	hub    repository.HubWebSocketRepository
	client repository.ClientWebSocketRepository
}

func NewWebsocketHandler(
	hub repository.HubWebSocketRepository,
	client repository.ClientWebSocketRepository,
) *WebsocketHandler {
	return &WebsocketHandler{
		hub:    hub,
		client: client,
	}
}

func (ws *WebsocketHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil) // conn is *websocket.Conn
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: clientの初期化

	go ws.client.WritePump()
	go ws.client.ReadPump()

	// TODO: hubにclientを登録
	// ws.hub.register <- client
}
