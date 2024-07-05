package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestMessageUseCase_ListMessages(t *testing.T) {
	t.Parallel()
	channelID := uuid.New().String()
	membershipID := uuid.New().String()
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockMembershipRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx       context.Context
			channelID string
			start     time.Time
			end       time.Time
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mcr.EXPECT().List(gomock.Any(), channelID, start, end).Return([]entity.Message{
					{
						ID:           "31894386-3e60-45a8-bc67-f46b72b42554",
						MembershipID: membershipID,
						Text:         "test message",
						CreatedAt:    time.Now(),
					},
				}, nil)
			},
			arg: struct {
				ctx       context.Context
				channelID string
				start     time.Time
				end       time.Time
			}{
				ctx:       context.Background(),
				channelID: channelID,
				start:     start,
				end:       end,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockMembershipRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			_, err := usecase.ListMessages(
				tt.arg.ctx,
				tt.arg.channelID,
				tt.arg.start,
				tt.arg.end,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageList() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_CreateMessage(t *testing.T) {
	t.Parallel()
	channelID := "f6bd2530-cd9b-4ac1-8dc1-38c697e6cce2"
	message := entity.Message{
		ID:   "31894386-3e60-45a8-bc67-f46b72b42554",
		Text: "test message",
	}

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockMembershipRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx       context.Context
			channelID string
			message   entity.Message
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mcr.EXPECT().Create(gomock.Any(), channelID, message).Return(nil)
			},
			arg: struct {
				ctx       context.Context
				channelID string
				message   entity.Message
			}{
				ctx:       context.Background(),
				channelID: channelID,
				message:   message,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockMembershipRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.CreateMessage(
				tt.arg.ctx,
				tt.arg.channelID,
				tt.arg.message,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_UpdateMessage(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	superMembershipID := uuid.New().String()
	notAuthorizedMembershipID := uuid.New().String()
	msgID := uuid.New().String()
	message := entity.Message{
		ID:           msgID,
		MembershipID: membershipID,
		Text:         "test message",
	}

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockMembershipRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx          context.Context
			message      entity.Message
			membershipID string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), membershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "test",
						IsAdmin:     false,
					}, nil)
				mcr.EXPECT().Update(gomock.Any(), message).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: membershipID,
			},
			wantErr: nil,
		},
		{
			name: "success: Super User",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), superMembershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "super_test",
						IsAdmin:     true,
					}, nil)
				mcr.EXPECT().Update(gomock.Any(), message).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: superMembershipID,
			},
			wantErr: nil,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), notAuthorizedMembershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "not_authorized_test",
						IsAdmin:     false,
					}, nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: notAuthorizedMembershipID,
			},
			wantErr: fmt.Errorf("don't have permission to update msg"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockMembershipRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.UpdateMessage(
				tt.arg.ctx,
				tt.arg.message,
				tt.arg.membershipID,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageUpdate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_DeleteMessage(t *testing.T) {
	t.Parallel()
	membershipID := uuid.New().String()
	superMembershipID := uuid.New().String()
	notAuthorizedMembershipID := uuid.New().String()
	channelID := uuid.New().String()
	msgID := uuid.New().String()
	message := entity.Message{
		ID:           msgID,
		MembershipID: membershipID,
		Text:         "test message",
	}

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockMembershipRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx          context.Context
			message      entity.Message
			membershipID string
			channelID    string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), membershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "test",
						IsAdmin:     false,
					}, nil)
				mcr.EXPECT().Delete(gomock.Any(), channelID, msgID).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
				channelID    string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: membershipID,
				channelID:    channelID,
			},
			wantErr: nil,
		},
		{
			name: "success: Super User",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), superMembershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "super_test",
						IsAdmin:     true,
					}, nil)
				mcr.EXPECT().Delete(gomock.Any(), channelID, msgID).Return(nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
				channelID    string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: superMembershipID,
				channelID:    channelID,
			},
			wantErr: nil,
		},
		{
			name: "Fail: Not authorized to delete",
			setup: func(
				mur *mock.MockMembershipRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), notAuthorizedMembershipID).
					Return(&entity.Membership{
						UserID:      uuid.New().String(),
						WorkspaceID: uuid.New().String(),
						Name:        "not_authorized_test",
						IsAdmin:     false,
					}, nil)
			},
			arg: struct {
				ctx          context.Context
				message      entity.Message
				membershipID string
				channelID    string
			}{
				ctx:          context.Background(),
				message:      message,
				membershipID: notAuthorizedMembershipID,
				channelID:    channelID,
			},
			wantErr: fmt.Errorf("don't have permission to delete msg"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockMembershipRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.DeleteMessage(
				tt.arg.ctx,
				tt.arg.message,
				tt.arg.membershipID,
				tt.arg.channelID,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
