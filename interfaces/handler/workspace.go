package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tusmasoma/connectHub-backend/interfaces/ws"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type WorkspaceHandler interface {
	CreateWorkspace(w http.ResponseWriter, r *http.Request)
}

type workspaceHandler struct {
	hm  *ws.HubManager
	auc usecase.AuthUseCase
	ruc usecase.RoomUseCase
	psr repository.PubSubRepository
	mcr repository.MessageCacheRepository
}

func NewWorkspaceHandler(
	hm *ws.HubManager,
	auc usecase.AuthUseCase,
	ruc usecase.RoomUseCase,
	psr repository.PubSubRepository,
	mcr repository.MessageCacheRepository,
) WorkspaceHandler {
	return &workspaceHandler{
		hm:  hm,
		auc: auc,
		ruc: ruc,
		psr: psr,
		mcr: mcr,
	}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

func (wh *workspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := wh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	var requestBody CreateWorkspaceRequest
	if ok := isValidCreateWorkspaceRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid workspace create request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid workspace create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hub := ws.NewHub(
		requestBody.Name,
		wh.ruc,
		wh.psr,
		wh.mcr,
	)

	go hub.Run()

	wh.hm.Add(hub.ID, hub)

	// membershipID := user.ID + "_" + hub.ID
	// hub.CreateRoom(ctx, membershipID, "general", false)
	// hub.CreateRoom(ctx, membershipID, "random", false)

	log.Info("Workspace created", log.Fstring("name", requestBody.Name))
	w.WriteHeader(http.StatusOK)
}

func isValidCreateWorkspaceRequest(body io.ReadCloser, requestBody *CreateWorkspaceRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.Name == "" {
		log.Info("Missing required fields", log.Fstring("name", requestBody.Name))
		return false
	}
	return true
}
