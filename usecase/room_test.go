package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestRoomUseCase_CreateRoom(t *testing.T) {
	t.Parallel()

	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6cce2"
	room := entity.Room{
		ID:          "roomID",
		Name:        "test",
		Description: "test",
		Private:     false,
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRoomRepository,
			m1 *mock.MockUserRoomRepository,
			m2 *mock.MockTransactionRepository,
		)
		arg struct {
			ctx         context.Context
			userID      string
			workspaceID string
			room        entity.Room
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockRoomRepository, urr *mock.MockUserRoomRepository, tr *mock.MockTransactionRepository) {
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				rr.EXPECT().Create(gomock.Any(), room).Return(nil)
				urr.EXPECT().Create(gomock.Any(), entity.UserRoom{
					UserWorkspaceID: userID + "_" + workspaceID,
					RoomID:          room.ID,
				}).Return(nil)
			},
			arg: struct {
				ctx         context.Context
				userID      string
				workspaceID string
				room        entity.Room
			}{
				ctx:         context.Background(),
				userID:      userID,
				workspaceID: workspaceID,
				room:        room,
			},
			wantErr: nil,
		},
		{
			name: "failed to create room",
			setup: func(rr *mock.MockRoomRepository, urr *mock.MockUserRoomRepository, tr *mock.MockTransactionRepository) {
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				rr.EXPECT().Create(gomock.Any(), room).Return(fmt.Errorf("failed to create room"))
			},
			arg: struct {
				ctx         context.Context
				userID      string
				workspaceID string
				room        entity.Room
			}{
				ctx:         context.Background(),
				userID:      userID,
				workspaceID: workspaceID,
				room:        room,
			},
			wantErr: fmt.Errorf("failed to create room"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			rr := mock.NewMockRoomRepository(ctrl)
			urr := mock.NewMockUserRoomRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr, urr, tr)
			}

			usecase := NewRoomUseCase(rr, urr, tr)
			err := usecase.CreateRoom(tt.arg.ctx, tt.arg.userID, tt.arg.workspaceID, tt.arg.room)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoomUseCase_ListUserWorkspaceRooms(t *testing.T) {
	t.Parallel()

	workspaceID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6cce2"

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRoomRepository,
		)
		arg struct {
			ctx         context.Context
			userID      string
			workspaceID string
		}
		want    []entity.Room
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockRoomRepository) {
				rr.EXPECT().ListUserWorkspaceRooms(gomock.Any(), userID, workspaceID).Return(
					[]entity.Room{
						{
							ID:          "roomID",
							Name:        "test",
							Description: "test",
							Private:     false,
						},
					},
					nil,
				)
			},
			arg: struct {
				ctx         context.Context
				userID      string
				workspaceID string
			}{
				ctx:         context.Background(),
				userID:      userID,
				workspaceID: workspaceID,
			},
			want: []entity.Room{
				{
					ID:          "roomID",
					Name:        "test",
					Description: "test",
					Private:     false,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			rr := mock.NewMockRoomRepository(ctrl)
			urr := mock.NewMockUserRoomRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr)
			}

			usecase := NewRoomUseCase(rr, urr, tr)
			getRooms, err := usecase.ListUserWorkspaceRooms(tt.arg.ctx, tt.arg.userID, tt.arg.workspaceID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ListUserWorkspaceRooms() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ListUserWorkspaceRooms() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && len(getRooms) != len(tt.want) {
				t.Errorf("ListUserWorkspaceRooms() = %v, want %v", getRooms, tt.want)
			}
		})
	}
}
