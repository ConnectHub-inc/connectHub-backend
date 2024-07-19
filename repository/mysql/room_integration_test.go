package mysql

import (
	"context"
	"testing"

	"github.com/doug-martin/goqu/v9"
)

func Test_ChannelRepository(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	userID := "5fe0e23e-6b49-11ee-b686-0242c0a87001"
	workspaceID := "5fe0e237-6b49-11ee-b686-0242c0a87001"
	membershipID := userID + "_" + workspaceID

	repo := NewChannelRepository(db, &dialect)

	channels, err := repo.ListMembershipChannels(ctx, membershipID)
	ValidateErr(t, err, nil)
	if len(channels) != 1 {
		t.Errorf("Expected 2 memberships, got %d", len(channels))
	}
	if channels[0].ID != "5fe0e239-6b49-11ee-b686-0242c0a87001" {
		t.Errorf("Expected channel ID to be 5fe0e239-6b49-11ee-b686-0242c0a87001, got %s", channels[0].ID)
	}
}
