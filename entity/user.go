package entity

import "github.com/google/uuid"

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ProfileImageURL string `json:"profile_image_url"`
	IsAdmin         bool   `json:"is_admin"`
}

func NewUser(name, email, password, profileImageURL string, isAdmin bool) *User {
	// TODO: validation
	return &User{
		ID:              uuid.New().String(),
		Name:            name,
		Email:           email,
		Password:        password,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
	}
}
