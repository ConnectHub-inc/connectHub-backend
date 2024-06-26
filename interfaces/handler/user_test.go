package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/usecase"
	"github.com/tusmasoma/connectHub-backend/usecase/mock"
)

func TestUserHandler_GetUser(t *testing.T) {
	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6cce2"
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockRoomUseCase,
			m2 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockRoomUseCase, m2 *mock.MockAuthUseCase) {
				m2.EXPECT().GetUserFromContext(gomock.Any()).Return(
					&entity.User{
						ID:       userID,
						Name:     "test",
						Email:    "test@gmail.com",
						Password: "password123",
					},
					nil,
				)
				m1.EXPECT().ListUserWorkspaceRooms(gomock.Any(), userID, workspaceID).Return(
					[]entity.Room{
						{
							ID:          "roomID",
							Name:        "test",
							Description: "test",
							Private:     false,
						},
					}, nil)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/user/get/%s", workspaceID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, ruc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/user/get/{workspace_id}", handler.GetUser)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_ListWorkspaceUsers(t *testing.T) {
	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m.EXPECT().ListWorkspaceUsers(
					gomock.Any(),
					workspaceID,
				).Return(
					[]entity.User{
						{
							ID:              "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
							Name:            "test",
							Email:           "test@gmail.com",
							ProfileImageURL: "https://test.com",
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/workspaces/%s/users", workspaceID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/workspaces/{workspace_id}/users", handler.ListWorkspaceUsers)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_ListRoomUsers(t *testing.T) {
	channelID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m.EXPECT().ListRoomUsers(
					gomock.Any(),
					channelID,
				).Return(
					[]entity.User{
						{
							ID:              "f6db2530-cd9b-4ac1-8dc1-38c795e6cce2",
							Name:            "test",
							Email:           "test@gmail.com",
							ProfileImageURL: "https://test.com",
						},
					},
					nil,
				)
			},
			in: func() *http.Request {
				url := fmt.Sprintf("/api/rooms/%s/users", channelID)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/rooms/{channel_id}/users", handler.ListRoomUsers)
			r.ServeHTTP(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m.EXPECT().CreateUserAndGenerateToken(
					gomock.Any(),
					"test@gmail.com",
					"password123",
				).Return(
					"eyJhbGciOiJIUzI1NiIsI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo.SflKxwRJSMeKKF2QT4fwpMeJf36P",
					nil,
				)
			},
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()
			handler.CreateUser(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	user := entity.User{
		ID:              userID,
		Name:            "test",
		Email:           "test@gmail.com",
		Password:        "password123",
		ProfileImageURL: "https://test.com",
	}
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					&user,
					nil,
				)
				m.EXPECT().UpdateUser(
					gomock.Any(),
					&usecase.UpdateUserParams{
						ID:              userID,
						Name:            "updated_test",
						Email:           "test@gmail.com",
						ProfileImageURL: "https://test.com",
					},
					user,
				).Return(nil)
			},
			in: func() *http.Request {
				userUpdateReq := UpdateUserRequest{
					ID:              userID,
					Name:            "updated_test",
					Email:           "test@gmail.com",
					ProfileImageURL: "https://test.com",
				}
				reqBody, _ := json.Marshal(userUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/user/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m1.EXPECT().GetUserFromContext(gomock.Any()).Return(
					&user,
					nil,
				)
			},
			in: func() *http.Request {
				userUpdateReq := UpdateUserRequest{
					ID: userID,
				}
				reqBody, _ := json.Marshal(userUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/user/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()
			handler.UpdateUser(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
			m1 *mock.MockAuthUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase, m1 *mock.MockAuthUseCase) {
				m.EXPECT().LoginAndGenerateToken(
					gomock.Any(),
					"test@gmail.com",
					"password123",
				).Return(
					"eyJhbGciOiJIUzI1NiIsI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo.SflKxwRJSMeKKF2QT4fwpMeJf36P",
					nil,
				)
			},
			in: func() *http.Request {
				userLoginReq := LoginRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userLoginReq := LoginRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userLoginReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)
			ruc := mock.NewMockRoomUseCase(ctrl)
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, ruc, auc)
			recorder := httptest.NewRecorder()
			handler.Login(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}
