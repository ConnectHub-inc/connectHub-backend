package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

type CreateUserAndGenerateTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

type UpdateUserArg struct {
	ctx    context.Context
	params *UpdateUserParams
	user   entity.User
}

func TestUserUseCase_ListWorkspaceUsers(t *testing.T) {
	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	users := []entity.User{
		{
			ID:              "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
			Name:            "test",
			Email:           "test@gmail.com",
			ProfileImageURL: "https://test.com",
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
		)
		arg     string
		want    []entity.User
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				m.EXPECT().ListWorkspaceUsers(
					gomock.Any(),
					workspaceID,
				).Return(users, nil)
			},
			arg:  workspaceID,
			want: users,
		},
		{
			name: "Fail: failed to list workspace users",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				m.EXPECT().ListWorkspaceUsers(
					gomock.Any(),
					workspaceID,
				).Return(nil, fmt.Errorf("failed to list workspace users"))
			},
			arg:     workspaceID,
			want:    nil,
			wantErr: fmt.Errorf("failed to list workspace users"),
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			users, err := usecase.ListWorkspaceUsers(context.Background(), tt.arg)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ListWorkspaceUsers() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ListWorkspaceUsers() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && len(users) != len(tt.want) {
				t.Errorf("ListWorkspaceUsers() = %v, want %v", users, tt.want)
			}
		})
	}
}

func TestUserUseCase_CreateUserAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
		)
		arg     CreateUserAndGenerateTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../.certificate/private_key.pem")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Email", Value: "test@gmail.com"}},
				).Return([]entity.User{}, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				m1.EXPECT().SetUserSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Email", Value: "test@gmail.com"}},
				).Return([]entity.User{{Name: "test", Email: "test@gmail.com"}}, nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			jwt, err := usecase.CreateUserAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
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
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
		)
		arg     UpdateUserArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				m.EXPECT().Update(
					gomock.Any(),
					userID,
					entity.User{
						ID:              userID,
						Name:            "update_test",
						Email:           "test@gmail.com",
						Password:        "password123",
						ProfileImageURL: "https://test.com",
					},
				).Return(nil)
			},
			arg: UpdateUserArg{
				ctx: context.Background(),
				params: &UpdateUserParams{
					ID:              userID,
					Name:            "update_test",
					Email:           "test@gmail.com",
					ProfileImageURL: "https://test.com",
				},
				user: user,
			},
			wantErr: nil,
		},
		{
			name: "Fail: don't have permission to update user",
			arg: UpdateUserArg{
				ctx: context.Background(),
				params: &UpdateUserParams{
					ID:              "f6db2530-cd9b-4ac1-8dc1-38c795e6eec3",
					Name:            "update_test",
					Email:           "test@gmail.com",
					ProfileImageURL: "https://test.com",
				},
				user: user,
			},
			wantErr: fmt.Errorf("don't have permission to update user"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			err := usecase.UpdateUser(tt.arg.ctx, tt.arg.params, tt.arg.user)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_LoginAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
		)
		arg     CreateUserAndGenerateTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../.certificate/private_key.pem")
				passward, _ := auth.PasswordEncrypt("password123")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Email", Value: "test@gmail.com"}},
				).Return(
					[]entity.User{
						{
							ID:       "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
							Name:     "test",
							Email:    "test@gmail.com",
							Password: passward,
						},
					}, nil,
				)
				m1.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return("", nil)
				m1.EXPECT().SetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
					gomock.Any(),
				).Return(nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: already logged in",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				passward, _ := auth.PasswordEncrypt("password123")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Email", Value: "test@gmail.com"}},
				).Return(
					[]entity.User{
						{
							ID:       "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
							Name:     "test",
							Email:    "test@gmail.com",
							Password: passward,
						},
					}, nil,
				)
				m1.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return("session", nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user id in cache"),
		},
		{
			name: "Fail: invalid passward",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository) {
				passward, _ := auth.PasswordEncrypt("password456")
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{{Field: "Email", Value: "test@gmail.com"}},
				).Return(
					[]entity.User{
						{
							ID:       "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
							Name:     "test",
							Email:    "test@gmail.com",
							Password: passward,
						},
					}, nil,
				)
				m1.EXPECT().GetUserSession(
					gomock.Any(),
					"f6db2530-cd9b-4ac1-8dc1-38c795e6eec2",
				).Return("", nil)
			},
			arg: CreateUserAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr)
			jwt, err := usecase.LoginAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}
