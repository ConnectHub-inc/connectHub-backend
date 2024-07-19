package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/usecase"
	"github.com/tusmasoma/connectHub-backend/usecase/mock"
)

func TestMembershipHandler_GetMembership(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	userID := uuid.New().String()
	membershipID := userID + "_" + workspaceID
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipUseCase,
			m1 *mock.MockChannelUseCase,
			m2 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockChannelUseCase, m2 *mock.MockAuthUseCase) {
				m2.EXPECT().GetUserFromContext(gomock.Any()).Return(
					&entity.User{
						ID:       userID,
						Email:    "test@gmail.com",
						Password: "password123",
					},
					nil,
				)
				m.EXPECT().GetMembership(gomock.Any(), membershipID).Return(
					&entity.Membership{
						UserID:          userID,
						WorkspaceID:     workspaceID,
						Name:            "test",
						ProfileImageURL: "https://test.jpg",
					}, nil,
				)
				m1.EXPECT().ListMembershipChannels(gomock.Any(), membershipID).Return(
					[]entity.Channel{
						{
							ID:          "channelID",
							Name:        "test",
							Description: "test",
							Private:     false,
						},
					}, nil)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/membership/get/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			muc := mock.NewMockMembershipUseCase(ctrl)
			ruc := mock.NewMockChannelUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(muc, ruc, auc)
			}

			handler := NewMembershipHandler(muc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/membership/get/{workspace_id}", handler.GetMembership)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestMembershipHandler_ListMemberships(t *testing.T) {
	t.Parallel()
	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipUseCase) {
				m.EXPECT().ListMemberships(
					gomock.Any(),
					workspaceID,
				).Return(
					[]entity.Membership{
						{
							UserID:          uuid.New().String(),
							WorkspaceID:     workspaceID,
							Name:            "test",
							ProfileImageURL: "https://test.com",
							IsAdmin:         false,
							IsDeleted:       false,
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/membership/list/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			muc := mock.NewMockMembershipUseCase(ctrl)
			ruc := mock.NewMockChannelUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(muc)
			}

			handler := NewMembershipHandler(muc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/membership/list/{workspace_id}", handler.ListMemberships)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestMembershipHandler_ListChannelMemberships(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	channelID := uuid.New().String()
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockAuthUseCase) {
				m.EXPECT().ListChannelMemberships(
					gomock.Any(),
					channelID,
				).Return(
					[]entity.Membership{
						{
							UserID:          uuid.New().String(),
							WorkspaceID:     workspaceID,
							Name:            "test",
							ProfileImageURL: "https://test.com",
							IsAdmin:         false,
							IsDeleted:       false,
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/membership/list-channel/%s", channelID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			muc := mock.NewMockMembershipUseCase(ctrl)
			ruc := mock.NewMockChannelUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(muc, auc)
			}

			handler := NewMembershipHandler(muc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/membership/list-channel/{channel_id}", handler.ListChannelMemberships)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestMembershipHandler_CreateMembership(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	workspaceID := uuid.New().String()
	user := &entity.User{
		ID:       userID,
		Email:    "test@gmail.com",
		Password: "password123",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					user,
					nil,
				)
				m.EXPECT().CreateMembership(
					gomock.Any(),
					&usecase.CreateMembershipParams{
						UserID:          userID,
						WorkspaceID:     workspaceID,
						Name:            "test",
						ProfileImageURL: "https://test.com",
						IsAdmin:         false,
					},
				).Return(nil)
			},
			in: func() *http.Request {
				userCreateReq := CreateMembershipRequest{
					Name:            "test",
					ProfileImageURL: "https://test.com",
					IsAdmin:         false,
				}
				reqBody, _ := json.Marshal(userCreateReq)
				url := fmt.Sprintf("/api/membership/create/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					user,
					nil,
				)
			},
			in: func() *http.Request {
				userCreateReq := CreateMembershipRequest{
					Name: "",
				}
				reqBody, _ := json.Marshal(userCreateReq)
				url := fmt.Sprintf("/api/membership/create/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			muc := mock.NewMockMembershipUseCase(ctrl)
			ruc := mock.NewMockChannelUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(muc, auc)
			}

			handler := NewMembershipHandler(muc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/api/membership/create/{workspace_id}", handler.CreateMembership)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestMembershipHandler_UpdateMembership(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	workspaceID := uuid.New().String()
	membershipID := userID + "_" + workspaceID
	user := &entity.User{
		ID:       userID,
		Email:    "test@gmail.com",
		Password: "password123",
	}
	membership := &entity.Membership{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            "test",
		ProfileImageURL: "https://test.jpg",
	}
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					user,
					nil,
				)
				m.EXPECT().GetMembership(gomock.Any(), membershipID).Return(
					membership, nil,
				)
				m.EXPECT().UpdateMembership(
					gomock.Any(),
					&usecase.UpdateMembershipParams{
						UserID:          userID,
						Name:            "updated_test",
						ProfileImageURL: "https://test.com",
					},
					*membership,
				).Return(nil)
			},
			in: func() *http.Request {
				userUpdateReq := UpdateMembershipRequest{
					UserID:          userID,
					Name:            "updated_test",
					ProfileImageURL: "https://test.com",
				}
				reqBody, _ := json.Marshal(userUpdateReq)
				url := fmt.Sprintf("/api/membership/update/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *mock.MockMembershipUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					user,
					nil,
				)
				m.EXPECT().GetMembership(gomock.Any(), membershipID).Return(
					membership, nil,
				)
			},
			in: func() *http.Request {
				userUpdateReq := UpdateMembershipRequest{
					UserID: "",
				}
				reqBody, _ := json.Marshal(userUpdateReq)
				url := fmt.Sprintf("/api/membership/update/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			muc := mock.NewMockMembershipUseCase(ctrl)
			ruc := mock.NewMockChannelUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(muc, auc)
			}

			handler := NewMembershipHandler(muc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/membership/update/{workspace_id}", handler.UpdateMembership)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
