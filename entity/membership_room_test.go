package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewMembershipChannel(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			membershipID string
			channelID    string
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				membershipID string
				channelID    string
			}{
				membershipID: "1",
				channelID:    "1",
			},
			wantErr: nil,
		},
		{
			name: "Fail: membershipID is required",
			arg: struct {
				membershipID string
				channelID    string
			}{
				membershipID: "",
				channelID:    "1",
			},
			wantErr: fmt.Errorf("membershipID is required"),
		},
		{
			name: "Fail: channelID is required",
			arg: struct {
				membershipID string
				channelID    string
			}{
				membershipID: "1",
				channelID:    "",
			},
			wantErr: fmt.Errorf("channelID is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMembershipChannel(tt.arg.membershipID, tt.arg.channelID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewMembershipChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
