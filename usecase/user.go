//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type UserUseCase interface {
	CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error)
	LoginAndGenerateToken(ctx context.Context, email string, passward string) (string, error)
	LogoutUser(ctx context.Context, userID string) error
}

type userUseCase struct {
	ur repository.UserRepository
	cr repository.UserCacheRepository
}

func NewUserUseCase(ur repository.UserRepository, cr repository.UserCacheRepository) UserUseCase {
	return &userUseCase{
		ur: ur,
		cr: cr,
	}
}

func (uuc *userUseCase) CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error) {
	user, err := uuc.CreateUser(ctx, email, passward)
	if err != nil {
		log.Printf("Internal server error while creating user")
		return "", err
	}

	jwt, jti := auth.GenerateToken(user.ID, user.Email)
	if err = uuc.cr.SetUserSession(ctx, user.ID, jti); err != nil {
		log.Print("Failed to set access token in cache")
		return "", err
	}

	return jwt, nil
}

func (uuc *userUseCase) CreateUser(ctx context.Context, email string, passward string) (*entity.User, error) {
	users, err := uuc.ur.List(ctx, []repository.QueryCondition{{Field: "Email", Value: email}})
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return nil, err
	}
	if len(users) > 0 {
		log.Printf("User with this name already exists - status: %d", http.StatusConflict)
		return nil, fmt.Errorf("user with this email already exists")
	}

	var user entity.User
	user.Email = email
	user.Name = auth.ExtractUsernameFromEmail(email)
	password, err := auth.PasswordEncrypt(passward)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		return nil, err
	}
	user.Password = password

	if err = uuc.ur.Create(ctx, user); err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (uuc *userUseCase) LoginAndGenerateToken(ctx context.Context, email string, passward string) (string, error) {
	var user entity.User
	// emailでMySQLにユーザー情報問い合わせ
	users, err := uuc.ur.List(ctx, []repository.QueryCondition{{Field: "Email", Value: email}})
	if err != nil {
		log.Printf("Error retrieving user by email")
		return "", err
	}
	if len(users) > 0 {
		user = users[0]
	}
	// 既にログイン済みかどうか確認する
	session, _ := uuc.cr.GetUserSession(ctx, user.ID)
	if session != "" {
		log.Printf("Already logged in")
		return "", fmt.Errorf("user id in cache")
	}

	// Clientから送られてきたpasswordをハッシュ化したものとMySQLから返されたハッシュ化されたpasswordを比較する
	if err = auth.CompareHashAndPassword(user.Password, passward); err != nil {
		log.Printf("password does not match")
		return "", err
	}

	jwt, jti := auth.GenerateToken(user.ID, email)
	if err = uuc.cr.SetUserSession(ctx, user.ID, jti); err != nil {
		log.Print("Failed to set access token in cache")
		return "", err
	}
	return jwt, nil
}

func (uuc *userUseCase) LogoutUser(ctx context.Context, userID string) error {
	if err := uuc.cr.Delete(ctx, userID); err != nil {
		log.Panicf("Failed to delete userID from cache")
		return err
	}
	return nil
}
