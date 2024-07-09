package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
)

func Test_MembershipRepository(t *testing.T) {
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

	room := entity.Room{
		ID:          channelID,
		WorkspaceID: workspaceID,
		Name:        "test",
		Description: "test",
		Private:     false,
	}

	membershipRoom := entity.MembershipRoom{
		MembershipID: membershipID,
		RoomID:       channelID,
	}

	userRepo := NewUserRepository(db, &dialect)
	membershipRepo := NewMembershipRepository(db, &dialect)
	roomRepo := NewRoomRepository(db, &dialect)
	membershipRoomRepo := NewMembershipRoomRepository(db, &dialect)

	// test get
	err := userRepo.Create(ctx, user)
	ValidateErr(t, err, nil)

	err = membershipRepo.Create(ctx, membership)
	ValidateErr(t, err, nil)

	getMembership, err := membershipRepo.Get(ctx, membershipID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(*getMembership, membership) {
		t.Errorf("Get() got = %v, want %v", getMembership, membership)
	}

	// test update
	membership.Name = "updated"
	err = membershipRepo.Update(ctx, membership)
	ValidateErr(t, err, nil)

	getMembership, err = membershipRepo.Get(ctx, membershipID)
	ValidateErr(t, err, nil)
	if getMembership.Name != "updated" {
		t.Errorf("Expected membership name 'updated', got %s", getMembership.Name)
	}

	// test list room memberships
	err = roomRepo.Create(ctx, room)
	ValidateErr(t, err, nil)

	err = membershipRoomRepo.Create(ctx, membershipRoom)
	ValidateErr(t, err, nil)

	getMemberships, err := membershipRepo.ListRoomMemberships(ctx, channelID)
	ValidateErr(t, err, nil)
	if len(getMemberships) != 1 {
		t.Errorf("Expected 1 membership, got %d", len(getMemberships))
	}

	// test soft delete
	err = membershipRepo.SoftDelete(ctx, membershipID)
	ValidateErr(t, err, nil)

	getMembership, err = membershipRepo.Get(ctx, membershipID)
	ValidateErr(t, err, nil)
	if !getMembership.IsDeleted {
		t.Errorf("Expected membership to be deleted, got %v", getMembership.IsDeleted)
	}

	// clean up
	err = membershipRoomRepo.Delete(ctx, membershipID, channelID)
	ValidateErr(t, err, nil)
	err = roomRepo.Delete(ctx, channelID)
	ValidateErr(t, err, nil)
	err = membershipRepo.Delete(ctx, membershipID)
	ValidateErr(t, err, nil)
	err = userRepo.Delete(ctx, userID)
	ValidateErr(t, err, nil)
}
