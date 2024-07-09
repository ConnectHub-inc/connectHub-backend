package entity

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/internal/log"
)

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func NewUser(email, password string) (*User, error) {
	if email == "" {
		log.Warn("Email is required", log.Fstring("email", email))
		return nil, fmt.Errorf("email is required")
	}
	if password == "" {
		log.Warn("Password is required", log.Fstring("password", password))
		return nil, fmt.Errorf("password is required")
	}
	password, err := auth.PasswordEncrypt(password)
	if err != nil {
		log.Error("Failed to encrypt password")
		return nil, err
	}
	return &User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: password,
	}, nil
}
