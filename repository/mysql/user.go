package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type userRepository struct {
	*base[entity.User]
}

func NewUserRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.UserRepository {
	return &userRepository{
		base: newBase[entity.User](db, dialect, "Users"),
	}
}

func (ur *userRepository) LockUserByEmail(ctx context.Context, email string) (bool, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var user entity.User
	query := "SELECT * FROM Users WHERE email = ? LIMIT 1 FOR UPDATE"
	row := executor.QueryRowContext(ctx, query, email)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No user found with the provided email", log.Fstring("email", email))
			return false, nil
		}
		log.Error("Failed to scan row", log.Ferror(err))
		return false, err
	}
	return true, nil
}
