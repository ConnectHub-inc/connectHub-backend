package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestAuthUseCase_FetchUserFromContext(t *testing.T) {
	t.Parallel()
	passward, _ := auth.PasswordEncrypt("password123")
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
		)
		in   func() context.Context
		want struct {
			user *entity.User
			err  error
		}
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().Get(
					gomock.Any(),
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
				).Return(
					&entity.User{
						ID:       "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
						Name:     "test",
						Email:    "test@gmail.com",
						Password: passward,
					}, nil,
				)
			},
			in: func() context.Context {
				ctx := context.WithValue(
					context.Background(),
					config.ContextUserIDKey,
					"5c5323e9-c78f-4dac-94ef-d34ab5ea8fed")
				return ctx
			},
			want: struct {
				user *entity.User
				err  error
			}{
				user: &entity.User{
					ID:       "5c5323e9-c78f-4dac-94ef-d34ab5ea8fed",
					Name:     "test",
					Email:    "test@gmail.com",
					Password: passward,
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: func() context.Context {
				ctx := context.Background()
				return ctx
			},
			want: struct {
				user *entity.User
				err  error
			}{
				user: nil,
				err:  fmt.Errorf("user name not found in request context"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		ctrl := gomock.NewController(t)
		mockUserRepo := mock.NewMockUserRepository(ctrl)

		if tt.setup != nil {
			tt.setup(mockUserRepo)
		}

		usecase := NewAuthUseCase(mockUserRepo)
		user, err := usecase.GetUserFromContext(tt.in())

		if (err != nil) != (tt.want.err != nil) {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
			t.Errorf("FetchUserFromContext() error = %v, wantErr %v", err, tt.want.err)
		}

		if !reflect.DeepEqual(user, tt.want.user) {
			t.Errorf("FetchUserFromContext() got = %v, want %v", *user, tt.want.user)
		}
	}
}
