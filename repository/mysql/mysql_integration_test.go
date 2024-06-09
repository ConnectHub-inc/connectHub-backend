package mysql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewMySQLDB(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *sql.DB
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
				t.Setenv("MYSQL_USER", "root")
				t.Setenv("MYSQL_PASSWORD", "connecthub")
				t.Setenv("MYSQL_HOST", "localhost")
				t.Setenv("MYSQL_PORT", mysqlPort)
				t.Setenv("MYSQL_DB_NAME", "connecthubdb")
			},
			want: &sql.DB{},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, _ := NewMySQLDB()

			if tt.want != nil {
				assert.NotNil(t, got, "Client should not be nil")
				err := got.Ping()
				require.NoError(t, err, "Error should be nil")
			} else {
				assert.Nil(t, got, "Client should be nil due to missing environment variables")
			}
		})
	}
}
