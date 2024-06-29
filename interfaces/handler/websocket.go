package handler

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/interfaces/ws"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.BufferSize,
	WriteBufferSize: config.BufferSize,
}

type WebsocketHandler struct {
	auc usecase.AuthUseCase
	psr repository.PubSubRepository
	muc usecase.MessageUseCase
	uru usecase.UserRoomUseCase
}

func NewWebsocketHandler(
	auc usecase.AuthUseCase,
	psr repository.PubSubRepository,
	muc usecase.MessageUseCase,
	uru usecase.UserRoomUseCase,
) *WebsocketHandler {
	return &WebsocketHandler{
		auc: auc,
		psr: psr,
		muc: muc,
		uru: uru,
	}
}

func (wsh *WebsocketHandler) WebSocket(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := wsh.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // conn is *websocket.Conn
	if err != nil {
		log.Error("Failed to upgrade connection", log.Ferror(err))
		return
	}

	client := ws.NewClient(user.ID, user.Name, conn, hub, wsh.psr, wsh.muc, wsh.uru)

	go client.WritePump()
	go client.ReadPump()

	hub.Register <- client
}
