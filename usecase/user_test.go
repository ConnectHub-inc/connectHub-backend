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

type SignUpAndGenerateTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

func TestUserUseCase_SignUpAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
			m2 *mock.MockTransactionRepository,
		)
		arg     SignUpAndGenerateTokenArg
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository, m2 *mock.MockTransactionRepository) {
				t.Setenv("PRIVATE_KEY_PATH", "../.certificate/private_key.pem")
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				m1.EXPECT().SetUserSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			arg: SignUpAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockUserCacheRepository, m2 *mock.MockTransactionRepository) {
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: SignUpAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr, tr)
			}

			usecase := NewUserUseCase(ur, cr, tr)
			jwt, err := usecase.SignUpAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("SignUpAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("SignUpAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

type LoginAndGenerateTokenArg struct {
	ctx      context.Context
	email    string
	passward string
}

func TestUserUseCase_LoginAndGenerateToken(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockUserCacheRepository,
		)
		arg     LoginAndGenerateTokenArg
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
			arg: LoginAndGenerateTokenArg{
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
			arg: LoginAndGenerateTokenArg{
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
			arg: LoginAndGenerateTokenArg{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			cr := mock.NewMockUserCacheRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, cr)
			}

			usecase := NewUserUseCase(ur, cr, tr)
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
