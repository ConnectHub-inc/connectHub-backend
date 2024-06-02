package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/dig"

	"github.com/tusmasoma/connectHub-backend/config"
)

func BuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		func(
			serverConfig *config.ServerConfig,
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

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}

	return container, nil
}
