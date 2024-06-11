package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestMessageUseCase_CreateMessage(t *testing.T) {
	t.Parallel()
	message := entity.Message{
		ID:      "31894386-3e60-45a8-bc67-f46b72b42554",
		Content: entity.MessageContent{Text: "test message"},
	}

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockUserRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx     context.Context
			message entity.Message
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mcr.EXPECT().Set(gomock.Any(), message.ID, message).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message entity.Message
			}{
				ctx:     context.Background(),
				message: message,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.CreateMessage(
				tt.arg.ctx,
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
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	msgID := "31894386-3e60-45a8-bc67-f46b72b42554"
	message := entity.Message{
		ID: msgID,
		Content: entity.MessageContent{
			UserID: userID,
			Text:   "test message",
		},
	}

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockUserRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx     context.Context
			message entity.Message
			userID  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), userID).
					Return(&entity.User{
						ID:      userID,
						Name:    "test",
						IsAdmin: false,
					}, nil)
				mcr.EXPECT().Set(gomock.Any(), msgID, message).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message entity.Message
				userID  string
			}{
				ctx:     context.Background(),
				message: message,
				userID:  userID,
			},
			wantErr: nil,
		},
		{
			name: "success: Super User",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), "f6db2530-cd9b-4ac1-8dc1-38c795e61234").
					Return(&entity.User{
						ID:      "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
						Name:    "super_test",
						IsAdmin: true,
					}, nil)
				mcr.EXPECT().Set(gomock.Any(), msgID, message).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message entity.Message
				userID  string
			}{
				ctx:     context.Background(),
				message: message,
				userID:  "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Not authorized to update",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), "f6db2530-cd9b-4ac1-8dc1-38c795e61234").
					Return(&entity.User{
						ID:      "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
						Name:    "not_authorized_test",
						IsAdmin: false,
					}, nil)
			},
			arg: struct {
				ctx     context.Context
				message entity.Message
				userID  string
			}{
				ctx:     context.Background(),
				message: message,
				userID:  "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
			},
			wantErr: fmt.Errorf("don't have permission to update msg"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.UpdateMessage(
				tt.arg.ctx,
				tt.arg.message,
				tt.arg.userID,
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
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	msgID := "31894386-3e60-45a8-bc67-f46b72b42554"

	patterns := []struct {
		name  string
		setup func(
			mur *mock.MockUserRepository,
			mmr *mock.MockMessageRepository,
			mcr *mock.MockMessageCacheRepository,
		)
		arg struct {
			ctx     context.Context
			content entity.MessageContent
			userID  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), userID).
					Return(&entity.User{
						ID:      userID,
						Name:    "test",
						IsAdmin: false,
					}, nil)
				mcr.EXPECT().Delete(gomock.Any(), msgID).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				content entity.MessageContent
				userID  string
			}{
				ctx: context.Background(),
				content: entity.MessageContent{
					UserID:    userID,
					MessageID: msgID,
				},
				userID: userID,
			},
			wantErr: nil,
		},
		{
			name: "success: Super User",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), "f6db2530-cd9b-4ac1-8dc1-38c795e61234").
					Return(&entity.User{
						ID:      "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
						Name:    "super_test",
						IsAdmin: true,
					}, nil)
				mcr.EXPECT().Delete(gomock.Any(), msgID).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				content entity.MessageContent
				userID  string
			}{
				ctx: context.Background(),
				content: entity.MessageContent{
					UserID:    userID,
					MessageID: msgID,
				},
				userID: "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Not authorized to delete",
			setup: func(
				mur *mock.MockUserRepository,
				mmr *mock.MockMessageRepository,
				mcr *mock.MockMessageCacheRepository,
			) {
				mur.EXPECT().Get(gomock.Any(), "f6db2530-cd9b-4ac1-8dc1-38c795e61234").
					Return(&entity.User{
						ID:      "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
						Name:    "not_authorized_test",
						IsAdmin: false,
					}, nil)
			},
			arg: struct {
				ctx     context.Context
				content entity.MessageContent
				userID  string
			}{
				ctx: context.Background(),
				content: entity.MessageContent{
					UserID:    userID,
					MessageID: msgID,
				},
				userID: "f6db2530-cd9b-4ac1-8dc1-38c795e61234",
			},
			wantErr: fmt.Errorf("don't have permission to delete msg"),
		},
	}
	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			mr := mock.NewMockMessageRepository(ctrl)
			mcr := mock.NewMockMessageCacheRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, mcr)
			}

			usecase := NewMessageUseCase(ur, mr, mcr)

			err := usecase.DeleteMessage(
				tt.arg.ctx,
				tt.arg.content,
				tt.arg.userID,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
