package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestUserUseCase_CreateUserRoom(t *testing.T) {
	t.Parallel()
	userID := "userID"
	roomID := "roomID"

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRoomRepository,
		)
		arg struct {
			ctx    context.Context
			userID string
			roomID string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockUserRoomRepository) {
				urr.EXPECT().Create(gomock.Any(),
					entity.UserRoom{
						UserID: userID,
						RoomID: roomID,
					},
				).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				userID string
				roomID string
			}{
				ctx:    context.Background(),
				userID: userID,
				roomID: roomID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockUserRoomRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewUserRoomUseCase(urr)

			err := usecase.CreateUserRoom(tt.arg.ctx, tt.arg.userID, tt.arg.roomID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_DeleteUserRoom(t *testing.T) {
	t.Parallel()
	userID := "userID"
	roomID := "roomID"

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRoomRepository,
		)
		arg struct {
			ctx    context.Context
			userID string
			roomID string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockUserRoomRepository) {
				urr.EXPECT().Delete(gomock.Any(), userID, roomID).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				userID string
				roomID string
			}{
				ctx:    context.Background(),
				userID: userID,
				roomID: roomID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockUserRoomRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewUserRoomUseCase(urr)

			err := usecase.DeleteUserRoom(tt.arg.ctx, tt.arg.userID, tt.arg.roomID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DeleteUserRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteUserRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
