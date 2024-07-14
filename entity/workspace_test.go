package entity

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestEntity_NewWorkspace(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			id   string
			name string
		}
		wantErr error
	}{
		{
			name: "success",
			arg: struct {
				id   string
				name string
			}{
				id:   uuid.New().String(),
				name: "test",
			},
			wantErr: nil,
		},
		{
			name: "Fail: id is required",
			arg: struct {
				id   string
				name string
			}{
				id:   "",
				name: "test",
			},
			wantErr: fmt.Errorf("id is required"),
		},
		{
			name: "Fail: name is required",
			arg: struct {
				id   string
				name string
			}{
				id:   uuid.New().String(),
				name: "",
			},
			wantErr: fmt.Errorf("name is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewWorkspace(tt.arg.id, tt.arg.name)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
