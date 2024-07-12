package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.uber.org/dig"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/interfaces/handler"
	"github.com/tusmasoma/connectHub-backend/interfaces/middleware"
	"github.com/tusmasoma/connectHub-backend/interfaces/ws"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository/mysql"
	"github.com/tusmasoma/connectHub-backend/repository/redis"
	"github.com/tusmasoma/connectHub-backend/usecase"
)

func BuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		log.Error("Failed to provide context")
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		config.NewCacheConfig,
		config.NewDBConfig,
		provideMySQLDialect,
		mysql.NewMySQLDB,
		mysql.NewTransactionRepository,
		mysql.NewUserRepository,
		mysql.NewMembershipRepository,
		mysql.NewMessageRepository,
		mysql.NewRoomRepository,
		mysql.NewMembershipRoomRepository,
		redis.NewRedisClient,
		redis.NewUserRepository,
		redis.NewMessageRepository,
		redis.NewPubSubRepository,
		usecase.NewUserUseCase,
		usecase.NewMembershipUseCase,
		usecase.NewMessageUseCase,
		usecase.NewRoomUseCase,
		usecase.NewMembershipRoomUseCase,
		usecase.NewAuthUseCase,
		handler.NewWebsocketHandler,
		handler.NewUserHandler,
		handler.NewMembershipHandler,
		middleware.NewAuthMiddleware,
		ws.NewHub,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
			membershipHandler handler.MembershipHandler,
			userHandler handler.UserHandler,
			authMiddleware middleware.AuthMiddleware,
			hub *ws.Hub,
		) *chi.Mux {
			r := chi.NewRouter()
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:     []string{"https://*", "http://*"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposedHeaders:     []string{"Link", "Authorization"},
				AllowCredentials:   true,
				MaxAge:             serverConfig.PreflightCacheDurationSec,
				OptionsPassthrough: true,
			}))

			go hub.Run()

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
					wsHandler.WebSocket(hub, w, r)
				})
			})

			// r.Use(middleware.Logging)
			r.Route("/api", func(r chi.Router) {
				r.Route("/membership", func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Get("/list/{workspace_id}", membershipHandler.ListMemberships)
					r.Get("/list-room/{channel_id}", membershipHandler.ListRoomMemberships)
					r.Get("/get/{workspace_id}", membershipHandler.GetMembership)
					r.Post("/create/{workspace_id}", membershipHandler.CreateMembership)
					r.Put("/update/{workspace_id}", membershipHandler.UpdateMembership)
				})

				r.Route("/user", func(r chi.Router) {
					r.Post("/signup", userHandler.SignUp)
					r.Post("/login", userHandler.Login)
					r.Group(func(r chi.Router) {
						r.Use(authMiddleware.Authenticate)
						r.Get("/logout", userHandler.Logout)
					})
				})
			})

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			log.Critical("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
			return nil, err
		}
	}

	log.Info("Container built successfully")
	return container, nil
}

func provideMySQLDialect() *goqu.DialectWrapper {
	dialect := goqu.Dialect("mysql")
	return &dialect
}
