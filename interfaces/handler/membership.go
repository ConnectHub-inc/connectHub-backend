package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

type MembershipHandler interface {
	GetMembership(w http.ResponseWriter, r *http.Request)
	ListMemberships(w http.ResponseWriter, r *http.Request)
	ListRoomMemberships(w http.ResponseWriter, r *http.Request)
	UpdateMembership(w http.ResponseWriter, r *http.Request)
}

type membershipHandler struct {
	muc usecase.MembershipUseCase
	ruc usecase.RoomUseCase
	auc usecase.AuthUseCase
}

func NewMembershipHandler(muc usecase.MembershipUseCase, ruc usecase.RoomUseCase, auc usecase.AuthUseCase) MembershipHandler {
	return &membershipHandler{
		muc: muc,
		ruc: ruc,
		auc: auc,
	}
}

type GetMembershipResponse struct {
	Name            string        `json:"name"`
	Email           string        `json:"email"`
	ProfileImageURL string        `json:"profile_image_url"`
	Rooms           []entity.Room `json:"rooms"`
}

func (mh *membershipHandler) GetMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspaceID := chi.URLParam(r, "workspace_id")
	user, err := mh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	membershipID := user.ID + "_" + workspaceID
	membership, err := mh.muc.GetMembership(ctx, membershipID)
	if err != nil {
		log.Error("Failed to get membership", log.Fstring("membershipID", membershipID))
		http.Error(w, "Failed to get membership", http.StatusInternalServerError)
		return
	}
	rooms, err := mh.ruc.ListMembershipRooms(ctx, membershipID)
	if err != nil {
		log.Error("Failed to list membership rooms", log.Fstring("membershipID", membershipID))
		http.Error(w, "Failed to list membership rooms", http.StatusInternalServerError)
		return
	}

	response := GetMembershipResponse{
		Name:            membership.Name,
		Email:           user.Email,
		ProfileImageURL: membership.ProfileImageURL,
		Rooms:           rooms,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Failed to encode membership to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode membership to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ListMembershipsResponse struct {
	Memberships []entity.Membership `json:"memberships"`
}

func (mh *membershipHandler) ListMemberships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workspaceID := chi.URLParam(r, "workspace_id")
	memberships, err := mh.muc.ListMemberships(ctx, workspaceID)
	if err != nil {
		log.Error("Failed to list memberships in workspace", log.Ferror(err))
		http.Error(w, "Failed to list memberships in workspace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ListMembershipsResponse{Memberships: memberships}); err != nil {
		log.Error("Failed to encode memberships to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode memberships to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ListRoomMembershipsResponse struct {
	Memberships []entity.Membership `json:"memberships"`
}

func (mh *membershipHandler) ListRoomMemberships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	channelID := chi.URLParam(r, "channel_id")
	memberships, err := mh.muc.ListRoomMemberships(ctx, channelID)
	if err != nil {
		log.Error("Failed to list room memberships", log.Fstring("channelID", channelID), log.Ferror(err))
		http.Error(w, "Failed to list room memberships", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ListRoomMembershipsResponse{Memberships: memberships}); err != nil {
		log.Error("Failed to encode memberships to JSON", log.Ferror(err))
		http.Error(w, "Failed to encode memberships to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type UpdateMembershipRequest struct {
	UserID          string `json:"userID"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (mh *membershipHandler) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspaceID := chi.URLParam(r, "workspace_id")
	user, err := mh.auc.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get UserInfo from context", log.Ferror(err))
		http.Error(w, fmt.Sprintf("Failed to get UserInfo from context: %v", err), http.StatusInternalServerError)
		return
	}

	membershipID := user.ID + "_" + workspaceID
	membership, err := mh.muc.GetMembership(ctx, membershipID)
	if err != nil {
		log.Error("Failed to get membership", log.Fstring("membershipID", membershipID))
		http.Error(w, "Failed to get membership", http.StatusInternalServerError)
		return
	}

	var requestBody UpdateMembershipRequest
	if ok := isValidUpdateMembershipRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid user udpate request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid user update request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := convertUpdateMembershipReqeuestToParams(requestBody)
	if err = mh.muc.UpdateMembership(ctx, params, *membership); err != nil {
		log.Error("Failed to update user", log.Ferror(err))
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidUpdateMembershipRequest(body io.ReadCloser, requestBody *UpdateMembershipRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.UserID == "" ||
		requestBody.Name == "" ||
		requestBody.ProfileImageURL == "" {
		log.Info(
			"Missing required fields",
			log.Fstring("userID", requestBody.UserID),
			log.Fstring("name", requestBody.Name),
			log.Fstring("profile_image_url", requestBody.ProfileImageURL),
		)
		return false
	}
	return true
}

func convertUpdateMembershipReqeuestToParams(req UpdateMembershipRequest) *usecase.UpdateMembershipParams {
	return &usecase.UpdateMembershipParams{
		UserID:          req.UserID,
		Name:            req.Name,
		ProfileImageURL: req.ProfileImageURL,
	}
}