package entity

import (
	"fmt"
	"testing"
)

func TestEntity_NewMessage(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  struct {
			membershipID string
			text         string
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				membershipID string
				text         string
			}{
				membershipID: "1",
				text:         "Hello, World!",
			},
			wantErr: nil,
		},
		{
			name: "Fail: membershipID is required",
			arg: struct {
				membershipID string
				text         string
			}{
				membershipID: "",
				text:         "Hello, World!",
			},
			wantErr: fmt.Errorf("id, membershipID and text are required"),
		},
		{
			name: "Fail: text is required",
			arg: struct {
				membershipID string
				text         string
			}{
				membershipID: "1",
				text:         "",
			},
			wantErr: fmt.Errorf("id, membershipID and text are required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMessage(tt.arg.membershipID, tt.arg.text)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntity_NewWSMessage(t *testing.T) {
	t.Parallel()
	message := Message{
		ID:           "1",
		MembershipID: "1",
		Text:         "Hello, World!",
	}

	patterns := []struct {
		name string
		arg  struct {
			action   string
			content  Message
			targetID string
			senderID string
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				action   string
				content  Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				content:  message,
				targetID: "1",
				senderID: "1",
			},
			wantErr: nil,
		},
		{
			name: "Fail: action is required",
			arg: struct {
				action   string
				content  Message
				targetID string
				senderID string
			}{
				action:   "",
				content:  message,
				targetID: "1",
				senderID: "1",
			},
			wantErr: fmt.Errorf("invalid action: %s", ""),
		},
		{
			name: "Fail: targetID is required",
			arg: struct {
				action   string
				content  Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				content:  message,
				targetID: "",
				senderID: "1",
			},
			wantErr: fmt.Errorf("action, targetID and senderID are required"),
		},
		{
			name: "Fail: senderID is required",
			arg: struct {
				action   string
				content  Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				content:  message,
				targetID: "1",
				senderID: "",
			},
			wantErr: fmt.Errorf("action, targetID and senderID are required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewWSMessage(tt.arg.action, tt.arg.content, tt.arg.targetID, tt.arg.senderID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewWSMessage() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewWSMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntity_NewWSMessages(t *testing.T) {
	t.Parallel()
	messages := []Message{
		{
			ID:           "1",
			MembershipID: "1",
			Text:         "Hello, World!",
		},
	}
	patterns := []struct {
		name string
		arg  struct {
			action   string
			contents []Message
			targetID string
			senderID string
		}
		wantErr error
	}{
		{
			name: "Success",
			arg: struct {
				action   string
				contents []Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				contents: messages,
				targetID: "1",
				senderID: "1",
			},
			wantErr: nil,
		},
		{
			name: "Fail: action is required",
			arg: struct {
				action   string
				contents []Message
				targetID string
				senderID string
			}{
				action:   "",
				contents: messages,
				targetID: "1",
				senderID: "1",
			},
			wantErr: fmt.Errorf("invalid action: %s", ""),
		},
		{
			name: "Fail: contents is required",
			arg: struct {
				action   string
				contents []Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				contents: nil,
				targetID: "1",
				senderID: "1",
			},
			wantErr: fmt.Errorf("action, contents, targetID and senderID are required"),
		},
		{
			name: "Fail: targetID is required",
			arg: struct {
				action   string
				contents []Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				contents: messages,
				targetID: "",
				senderID: "1",
			},
			wantErr: fmt.Errorf("action, contents, targetID and senderID are required"),
		},
		{
			name: "Fail: senderID is required",
			arg: struct {
				action   string
				contents []Message
				targetID string
				senderID string
			}{
				action:   "CREATE_MESSAGE",
				contents: messages,
				targetID: "1",
				senderID: "",
			},
			wantErr: fmt.Errorf("action, contents, targetID and senderID are required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewWSMessages(tt.arg.action, tt.arg.contents, tt.arg.targetID, tt.arg.senderID)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewWSMessages() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewWSMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
