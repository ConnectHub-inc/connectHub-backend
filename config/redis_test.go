package config

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
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
				t.Setenv("REDIS_ADDR", "localhost:6379")
				t.Setenv("REDIS_PASSWORD", "mypassword")
				t.Setenv("REDIS_DB", "0")
			},
			want: redis.NewClient(
				&redis.Options{
					Addr:     "localhost:6379",
					Password: "mypassword",
					DB:       0,
				}),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got := NewClient()

			if tt.want != nil {
				assert.NotNil(t, got, "Client should not be nil")
				assert.Equal(t, tt.want.Options().Addr, got.Options().Addr)
				assert.Equal(t, tt.want.Options().Password, got.Options().Password)
				assert.Equal(t, tt.want.Options().DB, got.Options().DB)
			} else {
				assert.Nil(t, got, "Client should be nil due to missing environment variables")
			}
		})
	}
}
