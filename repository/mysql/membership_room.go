package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type membershipRoomRepository struct {
	*base[entity.MembershipRoom]
}

func NewMembershipRoomRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.MembershipRoomRepository {
	return &membershipRoomRepository{
		base: newBase[entity.MembershipRoom](db, dialect, "Membership_Rooms"),
	}
}

func (mrr *membershipRoomRepository) Get(ctx context.Context, membershipID, roomID string) (*entity.MembershipRoom, error) {
	executor := mrr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var membershipRoom entity.MembershipRoom
	query, _, err := mrr.dialect.Select("*").From(mrr.tableName).Where(
		goqu.C("membership_id").Eq(membershipID),
		goqu.C("room_id").Eq(roomID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	row := executor.QueryRowContext(ctx, query)
	if err = mrr.structScanRow(&membershipRoom, row); err != nil {
		return nil, err
	}
	return &membershipRoom, nil
}

func (mrr *membershipRoomRepository) Delete(ctx context.Context, membershipID, roomID string) error {
	executor := mrr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := mrr.dialect.Delete(mrr.tableName).Where(
		goqu.C("membership_id").Eq(membershipID),
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
