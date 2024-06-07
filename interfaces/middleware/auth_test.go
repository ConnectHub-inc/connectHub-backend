package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func dummyTestHandler(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(config.ContextUserIDKey)
	if userID, _ := userIDValue.(string); userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	t.Setenv("PRIVATE_KEY_PATH", "../../.certificate/private_key.pem")
	t.Setenv("PUBLIC_KEY_PATH", "../../.certificate/public_key.pem")

	userID := uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2")
	email := "test@gmail.com"

	jwt, jti := auth.GenerateToken(userID.String(), email)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserCacheRepository,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserCacheRepository) {
				m.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(
					jti,
					nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+jwt)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: No Auth Header",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Auth Header Format",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", jwt)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Token",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+"invalid Token")
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: User ID Not In Cache",
			setup: func(m *mock.MockUserCacheRepository) {
				m.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(
					"",
					ErrCacheMiss,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+jwt)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: jti in Cache != jti in Payload",
			setup: func(m *mock.MockUserCacheRepository) {
				m.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return(
					"invalid jti",
					nil,
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+jwt)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			ctrl := gomock.NewController(t)
			repo := mock.NewMockUserCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(repo)
			}

			am := NewAuthMiddleware(repo)

			handler := am.Authenticate(http.HandlerFunc(dummyTestHandler))

			recoder := httptest.NewRecorder()
			handler.ServeHTTP(recoder, tt.in())

			// ステータスコードの検証
			if status := recoder.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
