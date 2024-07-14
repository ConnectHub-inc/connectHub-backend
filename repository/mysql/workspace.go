package mysql

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type workspaceRepository struct {
	*base[entity.Workspace]
}

func NewWorkspaceRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.WorkspaceRepository {
	return &workspaceRepository{
		base: newBase[entity.Workspace](db, dialect, "Workspaces"),
	}
}
