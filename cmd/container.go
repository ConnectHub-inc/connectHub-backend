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
		mysql.NewMessageRepository,
		mysql.NewRoomRepository,
		mysql.NewUserRoomRepository,
		redis.NewRedisClient,
		redis.NewUserRepository,
		redis.NewMessageRepository,
		redis.NewPubSubRepository,
		usecase.NewUserUseCase,
		usecase.NewMessageUseCase,
		usecase.NewRoomUseCase,
		usecase.NewUserRoomUseCase,
		usecase.NewAuthUseCase,
		handler.NewWebsocketHandler,
		handler.NewUserHandler,
		middleware.NewAuthMiddleware,
		ws.NewHub,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
			userHandler handler.UserHandler,
			authMiddleware middleware.AuthMiddleware,
			hub *ws.Hub,
		) *chi.Mux {
			r := chi.NewRouter()
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposedHeaders:   []string{"Link", "Authorization"},
				AllowCredentials: false,
				MaxAge:           serverConfig.PreflightCacheDurationSec,
			}))

			r.Use(middleware.Logging)

			go hub.Run()

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
					wsHandler.WebSocket(hub, w, r)
				})
			})

			r.Route("/api", func(r chi.Router) {
				r.Route("/workspaces", func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Get("/{workspace_id}/users", userHandler.ListWorkspaceUsers)
				})

				r.Route("/rooms", func(r chi.Router) {
					r.Use(authMiddleware.Authenticate)
					r.Get("/{room_id}/users", userHandler.ListWorkspaceUsers)
				})

				r.Route("/user", func(r chi.Router) {
					r.Post("/create", userHandler.CreateUser)
					r.Post("/login", userHandler.Login)
					r.Group(func(r chi.Router) {
						r.Use(authMiddleware.Authenticate)
						r.Get("/get/{workspace_id}", userHandler.GetUser)
						r.Put("/update", userHandler.UpdateUser)
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

func BuildContainer2(ctx context.Context) (*dig.Container, error) {
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
		mysql.NewMessageRepository,
		mysql.NewRoomRepository,
		mysql.NewUserRoomRepository,
		redis.NewRedisClient,
		redis.NewUserRepository,
		redis.NewMessageRepository,
		redis.NewPubSubRepository,
		usecase.NewUserUseCase,
		usecase.NewMessageUseCase,
		usecase.NewRoomUseCase,
		usecase.NewUserRoomUseCase,
		usecase.NewAuthUseCase,
		handler.NewWebsocketHandler,
		handler.NewUserHandler,
		middleware.NewAuthMiddleware,
		ws.NewHub,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
			userHandler handler.UserHandler,
			authMiddleware middleware.AuthMiddleware,
			hub *ws.Hub,
		) *http.ServeMux {
			mux := http.NewServeMux()
			go hub.Run()

			mux.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					wsHandler.WebSocket(hub, w, r)
				})).ServeHTTP(w, r)
			}))
			mux.Handle("/api/user/create", http.HandlerFunc(userHandler.CreateUser))
			return mux
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
