package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	ListWorkspaceUsers(w http.ResponseWriter, r *http.Request)
	ListRoomUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uur usecase.UserUseCase
	auc usecase.AuthUseCase
}

func NewUserHandler(uur usecase.UserUseCase, auc usecase.AuthUseCase) UserHandler {
	return &userHandler{
		uur: uur,
		auc: auc,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ListRoomUsersResponse struct {
	Users []entity.User `json:"users"`
}

func (uh *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := uh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Error("Failed to encode user to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *userHandler) ListWorkspaceUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workspaceID := chi.URLParam(r, "workspace_id")
	users, err := uh.uur.ListWorkspaceUsers(ctx, workspaceID)
	if err != nil {
		log.Error("Failed to list workspace users", log.Ferror(err))
		http.Error(w, "Failed to list workspace users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		log.Error("Failed to encode users to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *userHandler) ListRoomUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := chi.URLParam(r, "channel_id")
	users, err := uh.uur.ListRoomUsers(ctx, channelID)
	if err != nil {
		log.Error("Failed to list room users", log.Fstring("channelID", channelID), log.Ferror(err))
		http.Error(w, "Failed to list room users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ListRoomUsersResponse{Users: users}); err != nil {
		log.Error("Failed to encode users to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateUserRequest
	if ok := isValidCreateUserRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid user create request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid user create request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uur.CreateUserAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to create user and generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	w.WriteHeader(http.StatusOK)
}

func isValidCreateUserRequest(body io.ReadCloser, requestBody *CreateUserRequest) bool {
	// リクエストボディのJSONを構造体にデコード
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Info("Missing required fields", log.Fstring("email", requestBody.Email))
		return false
	}
	return true
}

func (uh *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := uh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	var requestBody UpdateUserRequest
	if ok := isValidUpdateUserRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid user udpate request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid user update request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := convertUpdateUserReqeuestToParams(requestBody)
	if err = uh.uur.UpdateUser(ctx, params, *user); err != nil {
		log.Error("Failed to update user", log.Ferror(err))
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidUpdateUserRequest(body io.ReadCloser, requestBody *UpdateUserRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.Name == "" ||
		requestBody.Email == "" ||
		requestBody.ProfileImageURL == "" {
		log.Info("Missing required fields", log.Fstring("email", requestBody.Email))
		return false
	}
	return true
}

func convertUpdateUserReqeuestToParams(req UpdateUserRequest) *usecase.UpdateUserParams {
	return &usecase.UpdateUserParams{
		ID:              req.ID,
		Name:            req.Name,
		Email:           req.Email,
		ProfileImageURL: req.ProfileImageURL,
	}
}

func (uh *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody LoginRequest
	if ok := isValidLoginRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid user login request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid user login request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uur.LoginAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to login or generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		http.Error(w, "Failed to Login or generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	w.WriteHeader(http.StatusOK)
}

func isValidLoginRequest(body io.ReadCloser, requestBody *LoginRequest) bool {
	// リクエストボディのJSONを構造体にデコード
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Info("Missing required fields", log.Fstring("email", requestBody.Email))
		return false
	}
	return true
}

func (uh *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := uh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	if err = uh.uur.LogoutUser(ctx, user.ID); err != nil {
		log.Error("Failed to logout", log.Ferror(err))
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
