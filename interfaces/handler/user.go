package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
