package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewWorkspace(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			name string
		}
		wantErr error
	}{
		{
			name: "success",
			arg: struct {
				name string
			}{
				name: "test",
			},
			wantErr: nil,
		},
		{
			name: "Fail: name is required",
			arg: struct {
				name string
			}{
				name: "",
			},
			wantErr: fmt.Errorf("name is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewWorkspace(tt.arg.name)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
