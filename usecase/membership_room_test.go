package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestUserUseCase_CreateMembershipRoom(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	roomID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRoomRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
			roomID       string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockMembershipRoomRepository) {
				urr.EXPECT().Create(gomock.Any(),
					entity.MembershipRoom{
						MembershipID: membershipID,
						RoomID:       roomID,
					},
				).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				membershipID string
				roomID       string
			}{
				ctx:          context.Background(),
				membershipID: membershipID,
				roomID:       roomID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockMembershipRoomRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewMembershipRoomUseCase(urr)

			err := usecase.CreateMembershipRoom(tt.arg.ctx, tt.arg.membershipID, tt.arg.roomID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_DeleteMembershipRoom(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	roomID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRoomRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
			roomID       string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockMembershipRoomRepository) {
				urr.EXPECT().Delete(gomock.Any(), membershipID, roomID).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				membershipID string
				roomID       string
			}{
				ctx:          context.Background(),
				membershipID: membershipID,
				roomID:       roomID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockMembershipRoomRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewMembershipRoomUseCase(urr)

			err := usecase.DeleteMembershipRoom(tt.arg.ctx, tt.arg.membershipID, tt.arg.roomID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DeleteMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
