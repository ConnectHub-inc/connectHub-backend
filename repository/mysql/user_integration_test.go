package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
)

func Test_UserRepository(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	repo := NewUserRepository(db, &dialect)

	user, err := entity.NewUser("test@gmail.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test CreateUser
	err = repo.Create(ctx, *user)
	ValidateErr(t, err, nil)

	// Test LockUserByEmail
	getUser, err := repo.LockUserByEmail(ctx, "test@gmail.com")
	ValidateErr(t, err, nil)
	if getUser == nil {
		t.Fatalf("Failed to get user by email")
	}

	// Test LockUserByEmail
	getUser, err = repo.LockUserByEmail(ctx, "fail@gmail.com")
	ValidateErr(t, err, sql.ErrNoRows)
	if getUser != nil {
		t.Fatalf("Failed to get user by email")
	}

	// clean up
	err = repo.Delete(ctx, user.ID)
	ValidateErr(t, err, nil)
}
