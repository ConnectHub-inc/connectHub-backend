package main

import (
	"context"
	"fmt"

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

func BuildContainer(ctx context.Context) (*dig.Container, error) { //nolint:funlen
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
		mysql.NewWorkspaceRepository,
		mysql.NewMessageRepository,
		mysql.NewChannelRepository,
		mysql.NewMembershipChannelRepository,
		redis.NewRedisClient,
		redis.NewUserRepository,
		redis.NewMessageRepository,
		redis.NewPubSubRepository,
		usecase.NewUserUseCase,
		usecase.NewMembershipUseCase,
		usecase.NewWorkspaceUseCase,
		usecase.NewMessageUseCase,
		usecase.NewChannelUseCase,
		usecase.NewMembershipChannelUseCase,
		usecase.NewAuthUseCase,
		ws.NewHubManager,
		handler.NewWebsocketHandler,
		handler.NewWorkspaceHandler,
		handler.NewUserHandler,
		handler.NewMembershipHandler,
		middleware.NewAuthMiddleware,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
			workspaceHandler handler.WorkspaceHandler,
			membershipHandler handler.MembershipHandler,
			userHandler handler.UserHandler,
			authMiddleware middleware.AuthMiddleware,
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

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/ws/{workspace_id}", wsHandler.WebSocket)
			})

			// r.Use(middleware.Logging)
			r.Route("/api", func(r chi.Router) {
				r.Route("/workspace", func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Post("/create", workspaceHandler.CreateWorkspace)
				})

				r.Route("/membership", func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Get("/list/{workspace_id}", membershipHandler.ListMemberships)
					r.Get("/list-channel/{channel_id}", membershipHandler.ListChannelMemberships)
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
