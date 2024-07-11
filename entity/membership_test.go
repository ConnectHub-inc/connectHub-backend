package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewMembership(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			userID          string
			workspaceID     string
			name            string
			profileImageURL string
			isAdmin         bool
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				userID          string
				workspaceID     string
				name            string
				profileImageURL string
				isAdmin         bool
			}{
				userID:          "1",
				workspaceID:     "1",
				name:            "test",
				profileImageURL: "https://www.hoge.com/avatar.jpg",
				isAdmin:         false,
			},
			wantErr: nil,
		},
		{
			name: "Success: profileImageURL is empty",
			arg: struct {
				userID          string
				workspaceID     string
				name            string
				profileImageURL string
				isAdmin         bool
			}{
				userID:          "1",
				workspaceID:     "1",
				name:            "test",
				profileImageURL: "",
				isAdmin:         false,
			},
			wantErr: nil,
		},
		{
			name: "Fail: userID is required",
			arg: struct {
				userID          string
				workspaceID     string
				name            string
				profileImageURL string
				isAdmin         bool
			}{
				userID:          "",
				workspaceID:     "1",
				name:            "test",
				profileImageURL: "https://www.hoge.com/avatar.jpg",
				isAdmin:         false,
			},
			wantErr: fmt.Errorf("userID is required"),
		},
		{
			name: "Fail: workspaceID is required",
			arg: struct {
				userID          string
				workspaceID     string
				name            string
				profileImageURL string
				isAdmin         bool
			}{
				userID:          "1",
				workspaceID:     "",
				name:            "test",
				profileImageURL: "https://www.hoge.com/avatar.jpg",
				isAdmin:         false,
			},
			wantErr: fmt.Errorf("workspaceID is required"),
		},
		{
			name: "Fail: name is required",
			arg: struct {
				userID          string
				workspaceID     string
				name            string
				profileImageURL string
				isAdmin         bool
			}{
				userID:          "1",
				workspaceID:     "1",
				name:            "",
				profileImageURL: "https://www.hoge.com/avatar.jpg",
				isAdmin:         false,
			},
			wantErr: fmt.Errorf("name is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMembership(tt.arg.userID, tt.arg.workspaceID, tt.arg.name, tt.arg.profileImageURL, tt.arg.isAdmin)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewMembership() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewMembership() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
