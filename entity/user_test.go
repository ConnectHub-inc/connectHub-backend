package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewUser(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			email    string
			password string
		}
		wantErr error
	}{
		{
			name: "success",
			arg: struct {
				email    string
				password string
			}{
				email:    "test@gmail.com",
				password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: email is required",
			arg: struct {
				email    string
				password string
			}{
				email:    "",
				password: "password123",
			},
			wantErr: fmt.Errorf("email is required"),
		},
		{
			name: "Fail: password is required",
			arg: struct {
				email    string
				password string
			}{
				email:    "test@gmail.com",
				password: "",
			},
			wantErr: fmt.Errorf("password is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewUser(tt.arg.email, tt.arg.password)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
