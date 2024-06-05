package mysql

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type roomRepository struct {
	*base[entity.Room]
}

func NewRoomRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.RoomRepository {
	return &roomRepository{
		base: newBase[entity.Room](db, dialect, "Rooms"),
	}
}
