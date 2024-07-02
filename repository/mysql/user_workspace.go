package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type userWorkspaceRepository struct {
	*base[entity.UserWorkspace]
}

func NewUserWorkspaceRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.UserWorkspaceRepository {
	return &userWorkspaceRepository{
		base: newBase[entity.UserWorkspace](db, dialect, "User_Workspaces"),
	}
}

func (uwr *userWorkspaceRepository) Get(ctx context.Context, userID, workspaceID string) (*entity.UserWorkspace, error) {
	executor := uwr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var userWorkspace entity.UserWorkspace
	query, _, err := uwr.dialect.Select("*").From(uwr.tableName).Where(
		goqu.C("user_id").Eq(userID),
		goqu.C("workspace_id").Eq(workspaceID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	row := executor.QueryRowContext(ctx, query)
	if err = uwr.structScanRow(&userWorkspace, row); err != nil {
		return nil, err
	}
	return &userWorkspace, nil
}
