package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestMembershipUseCase_ListMemberships(t *testing.T) {
	t.Parallel()

	workspaceID := uuid.New().String()
	memberships := []entity.Membership{
		{
			UserID:          uuid.New().String(),
			WorkspaceID:     workspaceID,
			Name:            "test",
			ProfileImageURL: "https://test.com",
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRepository,
		)
		arg struct {
			ctx         context.Context
			workspaceID string
		}
		want    []entity.Membership
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{
						{Field: "WorkspaceID", Value: workspaceID},
					},
				).Return(memberships, nil)
			},
			arg: struct {
				ctx         context.Context
				workspaceID string
			}{
				ctx:         context.Background(),
				workspaceID: workspaceID,
			},
			want:    memberships,
			wantErr: nil,
		},
		{
			name: "Fail: failed to list workspace memberships",
			setup: func(m *mock.MockMembershipRepository) {
				m.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{
						{Field: "WorkspaceID", Value: workspaceID},
					},
				).Return(nil, fmt.Errorf("failed to list workspace memberships"))
			},
			arg: struct {
				ctx         context.Context
				workspaceID string
			}{
				ctx:         context.Background(),
				workspaceID: workspaceID,
			},
			want:    nil,
			wantErr: fmt.Errorf("failed to list workspace memberships"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mr := mock.NewMockMembershipRepository(ctrl)
			mcr := mock.NewMockMembershipChannelRepository(ctrl)
			cr := mock.NewMockChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMembershipUseCase(mr, mcr, cr, tr)
			getMemberships, err := usecase.ListMemberships(tt.arg.ctx, tt.arg.workspaceID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ListMemberships() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ListMemberships() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && len(getMemberships) != len(tt.want) {
				t.Errorf("ListMemberships() = %v, want %v", getMemberships, tt.want)
			}
		})
	}
}

func TestMembershipUseCase_ListChannelMemberships(t *testing.T) {
	t.Parallel()

	channelID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	memberships := []entity.Membership{
		{
			UserID:          uuid.New().String(),
			WorkspaceID:     uuid.New().String(),
			Name:            "test",
			ProfileImageURL: "https://test.com",
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRepository,
		)
		arg struct {
			ctx       context.Context
			channelID string
		}
		want    []entity.Membership
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipRepository) {
				m.EXPECT().ListChannelMemberships(
					gomock.Any(),
					channelID,
				).Return(memberships, nil)
			},
			arg: struct {
				ctx       context.Context
				channelID string
			}{
				ctx:       context.Background(),
				channelID: channelID,
			},
			want:    memberships,
			wantErr: nil,
		},
		{
			name: "Fail: failed to list channel Memberships",
			setup: func(m *mock.MockMembershipRepository) {
				m.EXPECT().ListChannelMemberships(
					gomock.Any(),
					channelID,
				).Return(nil, fmt.Errorf("failed to list channel Memberships"))
			},
			arg: struct {
				ctx       context.Context
				channelID string
			}{
				ctx:       context.Background(),
				channelID: channelID,
			},
			want:    nil,
			wantErr: fmt.Errorf("failed to list channel Memberships"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mr := mock.NewMockMembershipRepository(ctrl)
			mcr := mock.NewMockMembershipChannelRepository(ctrl)
			cr := mock.NewMockChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMembershipUseCase(mr, mcr, cr, tr)
			getMemberships, err := usecase.ListChannelMemberships(tt.arg.ctx, tt.arg.channelID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ListChannelMemberships() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("ListChannelMemberships() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(getMemberships, tt.want) {
				t.Errorf("ListChannelMemberships() = %v, want %v", getMemberships, tt.want)
			}
		})
	}
}

func TestMembershipUseCase_CreateMembership(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	workspaceID := uuid.New().String()
	membershipID := userID + "_" + workspaceID
	channel1ID := uuid.New().String()
	channel2ID := uuid.New().String()
	membership := entity.Membership{
		ID:              membershipID,
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            "test",
		ProfileImageURL: "https://test.com",
		IsAdmin:         false,
	}
	channels := []entity.Channel{
		{
			ID:          channel1ID,
			WorkspaceID: workspaceID,
			Name:        "test1",
			Description: "test1",
			Private:     false,
		},
		{
			ID:          channel2ID,
			WorkspaceID: workspaceID,
			Name:        "test2",
			Description: "test2",
			Private:     false,
		},
	}
	membershipChannels := []entity.MembershipChannel{
		{
			MembershipID: membershipID,
			ChannelID:    channel1ID,
		},
		{
			MembershipID: membershipID,
			ChannelID:    channel2ID,
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRepository,
			m1 *mock.MockMembershipChannelRepository,
			m2 *mock.MockChannelRepository,
			m3 *mock.MockTransactionRepository,
		)
		arg struct {
			ctx    context.Context
			params *CreateMembershipParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockMembershipRepository,
				m1 *mock.MockMembershipChannelRepository,
				m2 *mock.MockChannelRepository,
				m3 *mock.MockTransactionRepository,
			) {
				m3.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().Create(
					gomock.Any(),
					membership,
				).Return(nil)
				m2.EXPECT().List(
					gomock.Any(),
					[]repository.QueryCondition{
						{
							Field: "workspace_id", Value: workspaceID,
						},
						{
							Field: "private", Value: false,
						},
					},
				).Return(channels, nil)
				m1.EXPECT().BatchCreate(
					gomock.Any(),
					membershipChannels,
				).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *CreateMembershipParams
			}{
				ctx: context.Background(),
				params: &CreateMembershipParams{
					UserID:          userID,
					WorkspaceID:     workspaceID,
					Name:            "test",
					ProfileImageURL: "https://test.com",
					IsAdmin:         false,
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
			mr := mock.NewMockMembershipRepository(ctrl)
			mcr := mock.NewMockMembershipChannelRepository(ctrl)
			cr := mock.NewMockChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr, mcr, cr, tr)
			}

			usecase := NewMembershipUseCase(mr, mcr, cr, tr)
			err := usecase.CreateMembership(tt.arg.ctx, tt.arg.params)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateMembership() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateMembership() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMembershipUseCase_UpdateMembership(t *testing.T) {
	t.Parallel()
	userID := uuid.New().String()
	notAuthorizedUserID := uuid.New().String()
	workspaceID := uuid.New().String()
	membership := entity.Membership{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            "test",
		ProfileImageURL: "https://test.com",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockMembershipRepository,
		)
		arg struct {
			ctx        context.Context
			params     *UpdateMembershipParams
			membership entity.Membership
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockMembershipRepository) {
				m.EXPECT().Update(
					gomock.Any(),
					entity.Membership{
						UserID:          userID,
						WorkspaceID:     workspaceID,
						Name:            "update_test",
						ProfileImageURL: "https://test.com",
					},
				).Return(nil)
			},
			arg: struct {
				ctx        context.Context
				params     *UpdateMembershipParams
				membership entity.Membership
			}{
				ctx: context.Background(),
				params: &UpdateMembershipParams{
					UserID:          userID,
					Name:            "update_test",
					ProfileImageURL: "https://test.com",
				},
				membership: membership,
			},
			wantErr: nil,
		},
		{
			name: "Fail: don't have permission to update Membership",
			arg: struct {
				ctx        context.Context
				params     *UpdateMembershipParams
				membership entity.Membership
			}{
				ctx: context.Background(),
				params: &UpdateMembershipParams{
					UserID:          notAuthorizedUserID,
					Name:            "update_test",
					ProfileImageURL: "https://test.com",
				},
				membership: membership,
			},
			wantErr: fmt.Errorf("don't have permission to update membership"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mr := mock.NewMockMembershipRepository(ctrl)
			mcr := mock.NewMockMembershipChannelRepository(ctrl)
			cr := mock.NewMockChannelRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMembershipUseCase(mr, mcr, cr, tr)
			err := usecase.UpdateMembership(tt.arg.ctx, tt.arg.params, tt.arg.membership)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UpdateMembership() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpdateMembership() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
