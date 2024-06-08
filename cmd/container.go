package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/dig"

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/interfaces/handler"
	"github.com/tusmasoma/connectHub-backend/interfaces/ws"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository/mysql"
	"github.com/tusmasoma/connectHub-backend/repository/redis"
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
		config.NewClient,
		config.NewDBConfig,
		config.NewDB,
		provideMySQLDialect,
		mysql.NewUserRepository,
		mysql.NewMessageRepository,
		mysql.NewRoomRepository,
		redis.NewUserRepository,
		redis.NewMessageRepository,
		redis.NewPubSubRepository,
		handler.NewWebsocketHandler,
		ws.NewHub,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
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

			go hub.Run()

			r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				wsHandler.WebSocket(hub, w, r)
			})

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			log.Fatal("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
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
