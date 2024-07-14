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
	wuc usecase.WorkspaceUseCase
	ruc usecase.RoomUseCase
	psr repository.PubSubRepository
	mcr repository.MessageCacheRepository
}

func NewWorkspaceHandler(
	hm *ws.HubManager,
	auc usecase.AuthUseCase,
	wuc usecase.WorkspaceUseCase,
	ruc usecase.RoomUseCase,
	psr repository.PubSubRepository,
	mcr repository.MessageCacheRepository,
) WorkspaceHandler {
	return &workspaceHandler{
		hm:  hm,
		auc: auc,
		wuc: wuc,
		ruc: ruc,
		psr: psr,
		mcr: mcr,
	}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type CreateWorkspaceResponse struct {
	ID   string `json:"workspace_id"`
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

	workspaceName := requestBody.Name
	hub := ws.NewHub(
		workspaceName,
		wh.ruc,
		wh.psr,
		wh.mcr,
	)

	go hub.Run()
	wh.hm.Add(hub.ID, hub)

	if err = wh.wuc.CreateWorkspace(ctx, hub.ID, workspaceName); err != nil {
		log.Error("Failed to create workspace", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to create workspace: %v", err), http.StatusInternalServerError)
		return
	}

	// membershipID := user.ID + "_" + hub.ID
	// hub.CreateRoom(ctx, membershipID, "general", false)
	// hub.CreateRoom(ctx, membershipID, "random", false)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(CreateWorkspaceResponse{ID: hub.ID, Name: workspaceName}); err != nil {
		log.Error("Failed to encode workspace to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode workspace to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("Successfully Workspace created", log.Fstring("workspaceID", hub.ID), log.Fstring("name", workspaceName))
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
