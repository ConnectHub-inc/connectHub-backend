package mysql

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type messageRepository struct {
	*base[entity.Message]
}

func NewSpotRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.MessageRepository {
	return &messageRepository{
		base: newBase[entity.Message](db, dialect, "Messages"),
	}
}
