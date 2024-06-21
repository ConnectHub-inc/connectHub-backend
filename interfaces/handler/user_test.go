package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/usecase"
	"github.com/tusmasoma/connectHub-backend/usecase/mock"
)

func TestUserHandler_GetUser(t *testing.T) {
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
					&entity.User{
						ID:       "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
						Name:     "test",
						Email:    "test@gmail.com",
						Password: "password123",
					},
					nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
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
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, auc)
			recorder := httptest.NewRecorder()
			handler.GetUser(recorder, tt.in())

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
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, auc)
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
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, auc)
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
			auc := mock.NewMockAuthUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc, auc)
			}

			handler := NewUserHandler(uuc, auc)
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
