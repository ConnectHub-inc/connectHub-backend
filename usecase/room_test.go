package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestChannelUseCase_CreateChannel(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	userID := uuid.New().String()
	membershipID := userID + "_" + workspaceID
	channelID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockChannelRepository,
			m1 *mock.MockMembershipChannelRepository,
			m2 *mock.MockTransactionRepository,
		)
		arg struct {
			ctx    context.Context
			params CreateChannelParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockChannelRepository, urr *mock.MockMembershipChannelRepository, tr *mock.MockTransactionRepository) {
				tr.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				rr.EXPECT().Create(
					gomock.Any(),
					entity.Channel{
						ID:          channelID,
						WorkspaceID: workspaceID,
						Name:        "test",
						Description: "test",
						Private:     false,
					},
				).Return(nil)
				urr.EXPECT().Create(gomock.Any(), entity.MembershipChannel{
					MembershipID: membershipID,
					ChannelID:    channelID,
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params CreateChannelParams
			}{
				ctx: context.Background(),
				params: CreateChannelParams{
					ID:           channelID,
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
			rr := mock.NewMockChannelRepository(ctrl)
			urr := mock.NewMockMembershipChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr, urr, tr)
			}

			usecase := NewChannelUseCase(rr, urr, tr)
			err := usecase.CreateChannel(tt.arg.ctx, tt.arg.params)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateChannel() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannelUseCase_ListMembershipChannels(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	userID := uuid.New().String()
	membershipID := userID + "_" + workspaceID

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockChannelRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
		}
		want    []entity.Channel
		wantErr error
	}{
		{
			name: "success",
			setup: func(rr *mock.MockChannelRepository) {
				rr.EXPECT().ListMembershipChannels(gomock.Any(), membershipID).Return(
					[]entity.Channel{
						{
							ID:          "channelID",
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
			want: []entity.Channel{
				{
					ID:          "channelID",
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
			rr := mock.NewMockChannelRepository(ctrl)
			urr := mock.NewMockMembershipChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(rr)
			}

			usecase := NewChannelUseCase(rr, urr, tr)
			getChannels, err := usecase.ListMembershipChannels(tt.arg.ctx, tt.arg.membershipID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ListUserWorkspaceChannels() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ListUserWorkspaceChannels() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && len(getChannels) != len(tt.want) {
				t.Errorf("ListUserWorkspaceChannels() = %v, want %v", getChannels, tt.want)
			}
		})
	}
}
