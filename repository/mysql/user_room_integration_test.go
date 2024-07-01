package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/entity"
)

func Test_UserRoomRepository(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	workspaceID := "5fe0e237-6b49-11ee-b686-0242c0a87001" // dml.test.sql
	userID := uuid.New().String()
	channelID := uuid.New().String()

	user := entity.User{
		ID:              userID,
		Name:            "test",
		Email:           "test@gmail.com",
		Password:        "test",
		ProfileImageURL: "https://test.com/test.jpg",
		IsAdmin:         false,
	}

	room := entity.Room{
		ID:          channelID,
		WorkspaceID: workspaceID,
		Name:        "test",
		Description: "test",
		Private:     false,
	}

	userRoom := entity.UserRoom{
		UserWorkspaceID: userID + "_" + workspaceID,
		RoomID:          channelID,
	}

	userRepo := NewUserRepository(db, &dialect)
	roomRepo := NewRoomRepository(db, &dialect)
	userRoomRepo := NewUserRoomRepository(db, &dialect)

	err := userRepo.Create(ctx, user)
	ValidateErr(t, err, nil)

	err = roomRepo.Create(ctx, room)
	ValidateErr(t, err, nil)

	err = userRoomRepo.Create(ctx, userRoom)
	ValidateErr(t, err, nil)

	getUserRoom, err := userRoomRepo.Get(ctx, userID, workspaceID, channelID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(*getUserRoom, userRoom) {
		t.Errorf("Get() = %v, want %v", getUserRoom, userRoom)
	}

	err = userRoomRepo.Delete(ctx, userID, workspaceID, channelID)
	ValidateErr(t, err, nil)

	getUserRoom, err = userRoomRepo.Get(ctx, userID, workspaceID, channelID)
	if err == nil {
		t.Errorf("Expected error for deleted item, got nil")
	}
	if getUserRoom != nil {
		t.Errorf("Expected nil for deleted item, got %v", getUserRoom)
	}
}
