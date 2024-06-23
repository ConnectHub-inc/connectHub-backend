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
	ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.User, error)
	ListRoomUsers(ctx context.Context, channelID string) ([]entity.User, error)
	CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams, user entity.User) error
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

func (uuc *userUseCase) ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.User, error) {
	users, err := uuc.ur.ListWorkspaceUsers(ctx, workspaceID)
	if err != nil {
		log.Error("Failed to list workspace users", log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return users, nil
}

func (uuc *userUseCase) ListRoomUsers(ctx context.Context, channelID string) ([]entity.User, error) {
	users, err := uuc.ur.ListRoomUsers(ctx, channelID)
	if err != nil {
		log.Error("Failed to list room users", log.Fstring("channelID", channelID))
		return nil, err
	}
	return users, nil
}

func (uuc *userUseCase) CreateUserAndGenerateToken(ctx context.Context, email string, passward string) (string, error) {
	user, err := uuc.CreateUser(ctx, email, passward)
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

func (uuc *userUseCase) CreateUser(ctx context.Context, email string, passward string) (*entity.User, error) {
	users, err := uuc.ur.List(ctx, []repository.QueryCondition{{Field: "Email", Value: email}})
	if err != nil {
		log.Error("Error retrieving user by email", log.Fstring("email", email))
		return nil, err
	}
	if len(users) > 0 {
		log.Info("User with this email already exists", log.Fstring("email", email))
		return nil, fmt.Errorf("user with this email already exists")
	}

	var user entity.User
	user.Email = email
	user.Name = auth.ExtractUsernameFromEmail(email)
	password, err := auth.PasswordEncrypt(passward)
	if err != nil {
		log.Error("Failed to encrypt password")
		return nil, err
	}
	user.Password = password

	if err = uuc.ur.Create(ctx, user); err != nil {
		log.Error("Failed to create user", log.Fstring("email", email))
		return nil, err
	}
	return &user, nil
}

type UpdateUserParams struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (uuc *userUseCase) UpdateUser(ctx context.Context, params *UpdateUserParams, user entity.User) error {
	if user.ID != params.ID {
		log.Warn("User don't have permission to update user", log.Fstring("userID", user.ID))
		return fmt.Errorf("don't have permission to update user")
	}

	user.Name = params.Name
	user.Email = params.Email
	user.ProfileImageURL = params.ProfileImageURL

	if err := uuc.ur.Update(ctx, user.ID, user); err != nil {
		log.Error("Failed to update user", log.Fstring("userID", user.ID))
		return err
	}
	return nil
}

func (uuc *userUseCase) LoginAndGenerateToken(ctx context.Context, email string, passward string) (string, error) {
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
	if err = auth.CompareHashAndPassword(user.Password, passward); err != nil {
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
