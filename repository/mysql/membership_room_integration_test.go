package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
)

func Test_MembershipChannelRepository(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	workspaceID := "5fe0e237-6b49-11ee-b686-0242c0a87001" // dml.test.sql
	userID := uuid.New().String()
	channelID := uuid.New().String()
	membershipID := userID + "_" + workspaceID

	user := entity.User{
		ID:       userID,
		Email:    "test@gmail.com",
		Password: "password123",
	}

	membership := entity.Membership{
		ID:              membershipID,
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            "test",
		ProfileImageURL: "https://test.com/test.jpg",
		IsAdmin:         false,
		IsDeleted:       false,
	}

	channel := entity.Channel{
		ID:          channelID,
		WorkspaceID: workspaceID,
		Name:        "test",
		Description: "test",
		Private:     false,
	}

	membershipChannel := entity.MembershipChannel{
		MembershipID: membershipID,
		ChannelID:    channelID,
	}

	userRepo := NewUserRepository(db, &dialect)
	membershipRepo := NewMembershipRepository(db, &dialect)
	channelRepo := NewChannelRepository(db, &dialect)
	membershipChannelRepo := NewMembershipChannelRepository(db, &dialect)

	err := userRepo.Create(ctx, user)
	ValidateErr(t, err, nil)

	err = membershipRepo.Create(ctx, membership)
	ValidateErr(t, err, nil)

	err = channelRepo.Create(ctx, channel)
	ValidateErr(t, err, nil)

	err = membershipChannelRepo.Create(ctx, membershipChannel)
	ValidateErr(t, err, nil)

	getMembershipChannel, err := membershipChannelRepo.Get(ctx, membershipID, channelID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(*getMembershipChannel, membershipChannel) {
		t.Errorf("Get() = %v, want %v", getMembershipChannel, membershipChannel)
	}

	err = membershipChannelRepo.Delete(ctx, membershipID, channelID)
	ValidateErr(t, err, nil)

	getMembershipChannel, err = membershipChannelRepo.Get(ctx, membershipID, channelID)
	if err == nil {
		t.Errorf("Expected error for deleted item, got nil")
	}
	if getMembershipChannel != nil {
		t.Errorf("Expected nil for deleted item, got %v", getMembershipChannel)
	}

	// clean up
	err = channelRepo.Delete(ctx, channelID)
	ValidateErr(t, err, nil)
	err = membershipRepo.Delete(ctx, membershipID)
	ValidateErr(t, err, nil)
	err = userRepo.Delete(ctx, userID)
	ValidateErr(t, err, nil)
}
