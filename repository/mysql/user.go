package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

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

func (ur *userRepository) LockUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := ur.dialect.Select("*").From(ur.tableName).Where(
		goqu.C("email").Eq(email),
	).ForUpdate(exp.NoWait).Limit(1).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	var user entity.User
	row := executor.QueryRowContext(ctx, query)
	if err = ur.structScanRow(&user, row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No user found with the provided email", log.Fstring("email", email))
			return nil, sql.ErrNoRows
		}
		log.Error("Failed to scan row", log.Ferror(err))
		return nil, err
	}
	return &user, nil
}
