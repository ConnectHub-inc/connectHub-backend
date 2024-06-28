package mysql

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type userRoomRepository struct {
	*base[entity.UserRoom]
}

func NewUserRoomRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.UserRoomRepository {
	return &userRoomRepository{
		base: newBase[entity.UserRoom](db, dialect, "User_Rooms"),
	}
}
