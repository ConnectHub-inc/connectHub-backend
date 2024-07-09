package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestRoomUseCase_CreateRoom(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	userID := uuid.New().String()
	membershipID := userID + "_" + workspaceID
	roomID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRoomRepository,
			m1 *mock.MockMembershipRoomRepository,
			m2 *mock.MockTransactionRepository,
		)
		arg struct {
			ctx    context.Context
			params CreateRoomParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockRoomRepository, urr *mock.MockMembershipRoomRepository, tr *mock.MockTransactionRepository) {
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				rr.EXPECT().Create(
					gomock.Any(),
					entity.Room{
						ID:          roomID,
						WorkspaceID: workspaceID,
						Name:        "test",
						Description: "test",
						Private:     false,
					},
				).Return(nil)
				urr.EXPECT().Create(gomock.Any(), entity.MembershipRoom{
					MembershipID: membershipID,
					RoomID:       roomID,
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params CreateRoomParams
			}{
				ctx: context.Background(),
				params: CreateRoomParams{
					ID:           roomID,
					MembershipID: membershipID,
					WorkspaceID:  workspaceID,
					Name:         "test",
					Description:  "test",
					Private:      false,
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
			urr := mock.NewMockMembershipRoomRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr, urr, tr)
			}

			usecase := NewRoomUseCase(rr, urr, tr)
			err := usecase.CreateRoom(tt.arg.ctx, tt.arg.params)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoomUseCase_ListMembershipRooms(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	userID := uuid.New().String()
	membershipID := userID + "_" + workspaceID

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockRoomRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
		}
		want    []entity.Room
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockRoomRepository) {
				rr.EXPECT().ListMembershipRooms(gomock.Any(), membershipID).Return(
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
				ctx          context.Context
				membershipID string
			}{
				ctx:          context.Background(),
				membershipID: membershipID,
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
			urr := mock.NewMockMembershipRoomRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr)
			}

			usecase := NewRoomUseCase(rr, urr, tr)
			getRooms, err := usecase.ListMembershipRooms(tt.arg.ctx, tt.arg.membershipID)

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
