package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
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

func (urr *userRoomRepository) Get(ctx context.Context, userID, workspaceID, roomID string) (*entity.UserRoom, error) {
	executor := urr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var userRoom entity.UserRoom
	userWorkspaceID := userID + "_" + workspaceID
	query, _, err := urr.dialect.Select("*").From(urr.tableName).Where(
		goqu.C("user_workspace_id").Eq(userWorkspaceID),
		goqu.C("room_id").Eq(roomID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	row := executor.QueryRowContext(ctx, query)
	if err = urr.structScanRow(&userRoom, row); err != nil {
		return nil, err
	}
	return &userRoom, nil
}

func (urr *userRoomRepository) Delete(ctx context.Context, userID, workspaceID, roomID string) error {
	executor := urr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	userWorkspaceID := userID + "_" + workspaceID
	query, _, err := urr.dialect.Delete(urr.tableName).Where(
		goqu.C("user_workspace_id").Eq(userWorkspaceID),
		goqu.C("room_id").Eq(roomID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return err
	}

	_, err = executor.ExecContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return err
	}
	return nil
}
