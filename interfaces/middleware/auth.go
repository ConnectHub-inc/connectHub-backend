package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/internal/auth"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

var ErrCacheMiss = errors.New("cache: key not found")

type AuthMiddleware interface {
	Authenticate(nextFunc http.Handler) http.Handler
}

type authMiddleware struct {
	rr repository.UserCacheRepository
}

func NewAuthMiddleware(rr repository.UserCacheRepository) AuthMiddleware {
	return &authMiddleware{
		rr: rr,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (am *authMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// リクエストヘッダにAuthorizationが存在するか確認
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Info("Authentication failed: missing Authorization header")
			http.Error(w, "Authentication failed: missing Authorization header", http.StatusUnauthorized)
			return
		}

		// "Bearer "から始まるか確認
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn("Authorization failed: header format must be Bearer {token}")
			http.Error(w, "Authorization failed: header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		jwt := parts[1]

		// アクセストークンの検証
		err := auth.ValidateAccessToken(jwt)
		if err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// JWTからペイロード取得
		var payload auth.Payload
		payload, err = auth.GetPayloadFromToken(jwt)
		if err != nil {
			log.Warn("Authentication failed: invalid access token", log.Ferror(err))
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// 該当のuserIdが存在するかキャッシュに問い合わせ
		jti, err := am.rr.GetUserSession(ctx, payload.UserID)
		if errors.Is(err, ErrCacheMiss) {
			log.Warn("Authentication failed: userId does not exist in cache", log.Fstring("userId", payload.UserID))
			http.Error(w, "Authentication failed: userId does not exist in cache", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Error(
				"Authentication failed: failed to get userId from cache",
				log.Fstring("userId", payload.UserID),
				log.Ferror(err),
			)
			http.Error(w, "Authentication failed: failed to get userId from cache", http.StatusUnauthorized)
			return
		}

		// Redisから取得したjtiとJWTのjtiを比較
		if payload.JTI != jti {
			log.Warn(
				"Authentication failed: jwt does not match",
				log.Fstring("jwtJTI", payload.JTI),
				log.Fstring("cacheJTI", jti),
			)
			http.Error(w, "Authentication failed: jwt does not match", http.StatusUnauthorized)
			return
		}

		// TODO:今後有効期限の確認も行う

		// コンテキストに userID を保存
		ctx = context.WithValue(ctx, config.ContextUserIDKey, payload.UserID)

		log.Info("Successfully Authentication", log.Fstring("userID", payload.UserID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
