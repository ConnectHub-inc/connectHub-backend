package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func Test_PasswordEncrypt(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name        string
		password    string
		expectError error
	}{
		{
			name:        "empty password",
			password:    "",
			expectError: nil,
		},
		{
			name:        "normal password",
			password:    "password123",
			expectError: nil,
		},
		{
			name:        "complex password",
			password:    "p@$$w0rd!23",
			expectError: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hash, err := PasswordEncrypt(tt.password)
			require.ErrorIs(t, err, tt.expectError)
			if err == nil && hash == "" {
				t.Errorf("Expected non-empty hash, got empty")
			}
		})
	}
}

func Test_CompareHashAndPassword(t *testing.T) {
	t.Parallel()

	validPassword := "validPassword"
	hashedPassword, _ := PasswordEncrypt(validPassword)

	patterns := []struct {
		name        string
		hash        string
		password    string
		expectError error
	}{
		{
			name:        "valid hash and password",
			hash:        hashedPassword,
			password:    validPassword,
			expectError: nil,
		},
		{
			name:        "invalid password",
			hash:        hashedPassword,
			password:    "invalidPassword",
			expectError: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := CompareHashAndPassword(tt.hash, tt.password)
			require.ErrorIs(t, err, tt.expectError)
		})
	}
}

func Test_ExtractUsernameFromEmail(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name         string
		email        string
		expectedName string
	}{
		{
			name:         "valid email",
			email:        "example@gmail.com",
			expectedName: "example",
		},
		{
			name:         "email without domain",
			email:        "example@",
			expectedName: "example",
		},
		{
			name:         "email without username",
			email:        "@gmail.com",
			expectedName: "",
		},
		{
			name:         "empty email",
			email:        "",
			expectedName: "",
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			username := ExtractUsernameFromEmail(tt.email)
			if username != tt.expectedName {
				t.Errorf("ExtractUsernameFromEmail() = %v, want %v", username, tt.expectedName)
			}
		})
	}
}
