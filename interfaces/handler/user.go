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
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uuc usecase.UserUseCase
	ruc usecase.RoomUseCase
	auc usecase.AuthUseCase
}

func NewUserHandler(uuc usecase.UserUseCase, ruc usecase.RoomUseCase, auc usecase.AuthUseCase) UserHandler {
	return &userHandler{
		uuc: uuc,
		ruc: ruc,
		auc: auc,
	}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody SignUpRequest
	if ok := isValidSignUpRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid sign up request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid sign up request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uuc.SignUpAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to create user and generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Info("User sign up successfully", log.Fstring("email", requestBody.Email))
	w.Header().Set("Authorization", "Bearer "+jwt)
	w.WriteHeader(http.StatusOK)
}

func isValidSignUpRequest(body io.ReadCloser, requestBody *SignUpRequest) bool {
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

	jwt, err := uh.uuc.LoginAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to login or generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		http.Error(w, "Failed to Login or generate token", http.StatusInternalServerError)
		return
	}

	log.Info("User login successfully", log.Fstring("email", requestBody.Email))
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

	if err = uh.uuc.LogoutUser(ctx, user.ID); err != nil {
		log.Error("Failed to logout", log.Ferror(err))
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	log.Info("User logout successfully", log.Fstring("userID", user.ID))
	w.WriteHeader(http.StatusOK)
}
