package mysql

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
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
