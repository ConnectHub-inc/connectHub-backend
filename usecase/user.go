//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type UserUseCase interface {
	CreateUserAndGenerateToken(ctx context.Context, email string, password string) (string, error)
	LoginAndGenerateToken(ctx context.Context, email string, password string) (string, error)
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

func (uuc *userUseCase) CreateUserAndGenerateToken(ctx context.Context, email string, password string) (string, error) {
	user, err := uuc.CreateUser(ctx, email, password)
	if err != nil {
		log.Error("Failed to create user", log.Fstring("email", email))
		return "", err
	}

	jwt, jti := auth.GenerateToken(user.ID, user.Email)
	if err = uuc.cr.SetUserSession(ctx, user.ID, jti); err != nil {
		log.Error("Failed to set access token in cache", log.Fstring("userID", user.ID), log.Fstring("jti", jti))
		return "", err
	}

	return jwt, nil
}

func (uuc *userUseCase) CreateUser(ctx context.Context, email string, password string) (*entity.User, error) {
	users, err := uuc.ur.List(ctx, []repository.QueryCondition{{Field: "Email", Value: email}})
	if err != nil {
		log.Error("Error retrieving user by email", log.Fstring("email", email))
		return nil, err
	}
	if len(users) > 0 {
		log.Info("User with this email already exists", log.Fstring("email", email))
		return nil, fmt.Errorf("user with this email already exists")
	}

	user, err := entity.NewUser(email, password)
	if err != nil {
		log.Error("Failed to create user", log.Ferror(err))
		return nil, err
	}

	if err = uuc.ur.Create(ctx, *user); err != nil {
		log.Error("Failed to create user", log.Fstring("email", email))
		return nil, err
	}
	return user, nil
}

func (uuc *userUseCase) LoginAndGenerateToken(ctx context.Context, email string, password string) (string, error) {
	var user entity.User
	// emailでMySQLにユーザー情報問い合わせ
	users, err := uuc.ur.List(ctx, []repository.QueryCondition{{Field: "Email", Value: email}})
	if err != nil {
		log.Error("Error retrieving user by email", log.Fstring("email", email))
		return "", err
	}
	if len(users) > 0 {
		user = users[0]
	}
	// 既にログイン済みかどうか確認する
	session, _ := uuc.cr.GetUserSession(ctx, user.ID)
	if session != "" {
		log.Info("Already logged in", log.Fstring("userID", user.ID))
		return "", fmt.Errorf("user id in cache")
	}

	// Clientから送られてきたpasswordをハッシュ化したものとMySQLから返されたハッシュ化されたpasswordを比較する
	if err = auth.CompareHashAndPassword(user.Password, password); err != nil {
		log.Info("Password does not match", log.Fstring("email", email))
		return "", err
	}

	jwt, jti := auth.GenerateToken(user.ID, email)
	if err = uuc.cr.SetUserSession(ctx, user.ID, jti); err != nil {
		log.Error("Failed to set access token in cache", log.Fstring("userID", user.ID), log.Fstring("jti", jti))
		return "", err
	}
	return jwt, nil
}

func (uuc *userUseCase) LogoutUser(ctx context.Context, userID string) error {
	if err := uuc.cr.Delete(ctx, userID); err != nil {
		log.Error("Failed to delete userID from cache", log.Fstring("userID", userID))
		return err
	}
	log.Info("Successfully logged out", log.Fstring("userID", userID))
	return nil
}
