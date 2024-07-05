package entity

import "github.com/google/uuid"

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func NewUser(email, password string) User {
	// TODO: validation
	return User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: password,
	}
}
