package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Transaction(t *testing.T) {
	dialect := goqu.Dialect("mysql")
	ctx := context.Background()
	userID := uuid.NewString()
	items := []Item{
		{ID: uuid.NewString(), UserID: userID, Text: "bar", Count: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "bat", Text: "baz", Count: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "qux", Text: "quux", Count: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	repo := newBase[Item](db, &dialect, "TestItems")

	patterns := []struct {
		name        string
		fn          func(ctx context.Context) error
		item        Item
		wantTXError error
		wantError   error
	}{
		{
			name: "default",
			fn: func(ctx context.Context) error {
				err := repo.Create(ctx, items[0])
				return err
			},
			item:        items[0],
			wantTXError: nil,
			wantError:   nil,
		},
		{
			name: "rollback",
			fn: func(ctx context.Context) error {
				err := repo.Create(ctx, items[1])
				if err != nil {
					return err
				}
				// force error to trigger rollback
				return fmt.Errorf("forced error to trigger rollback")
			},
			item:        items[1],
			wantTXError: fmt.Errorf("forced error to trigger rollback"),
			wantError:   sql.ErrNoRows,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			trepo := NewTransactionRepository(db)

			err := trepo.Transaction(ctx, tt.fn)
			ValidateErr(t, err, tt.wantTXError)

			_, err = repo.Get(ctx, tt.item.ID)
			ValidateErr(t, err, tt.wantError)
		})
	}
}

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
				t.Setenv("MYSQL_DB_NAME", "connecthubTestDB")
			},
			want: &sql.DB{},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			ctx := context.Background()
			got, _ := NewMySQLDB(ctx)

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
