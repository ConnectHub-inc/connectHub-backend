package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewChannel(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			id          string
			workspaceID string
			name        string
			description string
			private     bool
		}
		wantErr error
	}{
		{
			name: "success",
			arg: struct {
				id          string
				workspaceID string
				name        string
				description string
				private     bool
			}{
				id:          "1",
				workspaceID: "1",
				name:        "test",
				description: "test",
				private:     false,
			},
			wantErr: nil,
		},
		{
			name: "Fail: ID is required",
			arg: struct {
				id          string
				workspaceID string
				name        string
				description string
				private     bool
			}{
				id:          "",
				workspaceID: "1",
				name:        "test",
				description: "test",
				private:     false,
			},
			wantErr: fmt.Errorf("id is required"),
		},
		{
			name: "Fail: workspaceID is required",
			arg: struct {
				id          string
				workspaceID string
				name        string
				description string
				private     bool
			}{
				id:          "1",
				workspaceID: "",
				name:        "test",
				description: "test",
				private:     false,
			},
			wantErr: fmt.Errorf("workspaceID is required"),
		},
		{
			name: "Fail: name is required",
			arg: struct {
				id          string
				workspaceID string
				name        string
				description string
				private     bool
			}{
				id:          "1",
				workspaceID: "1",
				name:        "",
				description: "test",
				private:     false,
			},
			wantErr: fmt.Errorf("name is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewChannel(tt.arg.id, tt.arg.workspaceID, tt.arg.name, tt.arg.description, tt.arg.private)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewChannel() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
