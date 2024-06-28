package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestRoomUseCase_ListUserWorkspaceRooms(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
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
