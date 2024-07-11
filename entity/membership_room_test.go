package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewMembershipRoom(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			membershipID string
			roomID       string
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				membershipID string
				roomID       string
			}{
				membershipID: "1",
				roomID:       "1",
			},
			wantErr: nil,
		},
		{
			name: "Fail: membershipID is required",
			arg: struct {
				membershipID string
				roomID       string
			}{
				membershipID: "",
				roomID:       "1",
			},
			wantErr: fmt.Errorf("membershipID is required"),
		},
		{
			name: "Fail: roomID is required",
			arg: struct {
				membershipID string
				roomID       string
			}{
				membershipID: "1",
				roomID:       "",
			},
			wantErr: fmt.Errorf("roomID is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMembershipRoom(tt.arg.membershipID, tt.arg.roomID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewMembershipRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
