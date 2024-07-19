package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type membershipChannelRepository struct {
	*base[entity.MembershipChannel]
}

func NewMembershipChannelRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.MembershipChannelRepository {
	return &membershipChannelRepository{
		base: newBase[entity.MembershipChannel](db, dialect, "Membership_Channels"),
	}
}

func (mrr *membershipChannelRepository) Get(ctx context.Context, membershipID, channelID string) (*entity.MembershipChannel, error) {
	executor := mrr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var membershipChannel entity.MembershipChannel
	query, _, err := mrr.dialect.Select("*").From(mrr.tableName).Where(
		goqu.C("membership_id").Eq(membershipID),
		goqu.C("channel_id").Eq(channelID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	row := executor.QueryRowContext(ctx, query)
	if err = mrr.structScanRow(&membershipChannel, row); err != nil {
		return nil, err
	}
	return &membershipChannel, nil
}

func (mrr *membershipChannelRepository) Delete(ctx context.Context, membershipID, channelID string) error {
	executor := mrr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := mrr.dialect.Delete(mrr.tableName).Where(
		goqu.C("membership_id").Eq(membershipID),
		goqu.C("channel_id").Eq(channelID),
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
