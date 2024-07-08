package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewRedisClient(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *redis.Client
		err   error
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: nil,
		},
		{
			name: "set env",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("REDIS_ADDR", fmt.Sprintf("localhost:%s", redisPort))
				t.Setenv("REDIS_PASSWORD", "")
				t.Setenv("REDIS_DB", "0")
			},
			want: redis.NewClient(
				&redis.Options{
					Addr:     fmt.Sprintf("localhost:%s", redisPort),
					Password: "",
					DB:       0,
				}),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			ctx := context.Background()
			got := NewRedisClient(ctx)

			if tt.want != nil {
				assert.NotNil(t, got, "Client should not be nil")
				assert.Equal(t, tt.want.Options().Addr, got.Options().Addr)
				assert.Equal(t, tt.want.Options().Password, got.Options().Password)
				assert.Equal(t, tt.want.Options().DB, got.Options().DB)

				_, err := got.Ping(context.Background()).Result()
				require.NoError(t, err, "Error should be nil")
			} else {
				assert.Nil(t, got, "Client should be nil due to missing environment variables")
			}
		})
	}
}
