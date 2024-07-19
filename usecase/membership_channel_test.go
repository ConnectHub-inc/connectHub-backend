package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestUserUseCase_CreateMembershipChannel(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	channelID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipChannelRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
			channelID    string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockMembershipChannelRepository) {
				urr.EXPECT().Create(gomock.Any(),
					entity.MembershipChannel{
						MembershipID: membershipID,
						ChannelID:    channelID,
					},
				).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				membershipID string
				channelID    string
			}{
				ctx:          context.Background(),
				membershipID: membershipID,
				channelID:    channelID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockMembershipChannelRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewMembershipChannelUseCase(urr)

			err := usecase.CreateMembershipChannel(tt.arg.ctx, tt.arg.membershipID, tt.arg.channelID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_DeleteMembershipChannel(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	channelID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipChannelRepository,
		)
		arg struct {
			ctx          context.Context
			membershipID string
			channelID    string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(urr *mock.MockMembershipChannelRepository) {
				urr.EXPECT().Delete(gomock.Any(), membershipID, channelID).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				membershipID string
				channelID    string
			}{
				ctx:          context.Background(),
				membershipID: membershipID,
				channelID:    channelID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			urr := mock.NewMockMembershipChannelRepository(ctrl)

			if tt.setup != nil {
				tt.setup(urr)
			}

			usecase := NewMembershipChannelUseCase(urr)

			err := usecase.DeleteMembershipChannel(tt.arg.ctx, tt.arg.membershipID, tt.arg.channelID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DeleteMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
