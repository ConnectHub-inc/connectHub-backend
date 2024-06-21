package mysql

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
)

func Test_UserRepository(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	workspaceID := "5fe0e237-6b49-11ee-b686-0242c0a87001"

	repo := NewUserRepository(db, &dialect)

	users, err := repo.ListWorkspaceUsers(ctx, workspaceID)
	ValidateErr(t, err, nil)
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}
