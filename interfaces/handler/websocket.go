package handler

import (
	"net/http"

	"github.com/go-chi/chi"
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
	hm  *ws.HubManager
	auc usecase.AuthUseCase
	psr repository.PubSubRepository
	muc usecase.MessageUseCase
	uru usecase.MembershipChannelUseCase
}

func NewWebsocketHandler(
	hm *ws.HubManager,
	auc usecase.AuthUseCase,
	psr repository.PubSubRepository,
	muc usecase.MessageUseCase,
	uru usecase.MembershipChannelUseCase,
) *WebsocketHandler {
	return &WebsocketHandler{
		hm:  hm,
		auc: auc,
		psr: psr,
		muc: muc,
		uru: uru,
	}
}

func (wsh *WebsocketHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := wsh.auc.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Failed to get UserInfo from context", http.StatusInternalServerError)
		return
	}

	workspaceID := chi.URLParam(r, "workspace_id")
	hub, exists := wsh.hm.Get(workspaceID)
	if !exists {
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // conn is *websocket.Conn
	if err != nil {
		log.Error("Failed to upgrade connection", log.Ferror(err))
		return
	}

	client := ws.NewClient(user.ID, conn, hub, wsh.psr, wsh.muc, wsh.uru)

	go client.WritePump()
	go client.ReadPump()

	hub.Register <- client

	log.Info(
		"Successfully Client connected",
		log.Fstring("userID", user.ID),
		log.Fstring("workspaceID", workspaceID),
	)
}
