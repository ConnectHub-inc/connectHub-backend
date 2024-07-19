package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type membershipRepository struct {
	*base[entity.Membership]
}

func NewMembershipRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.MembershipRepository {
	return &membershipRepository{
		base: newBase[entity.Membership](db, dialect, "Memberships"),
	}
}

func (mr *membershipRepository) Get(ctx context.Context, id string) (*entity.Membership, error) {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var membership entity.Membership
	userID, workspaceID, err := membership.SplitMembershipID(id)
	if err != nil {
		log.Error("Failed to split membership ID", log.Ferror(err))
		return nil, err
	}
	query, _, err := mr.dialect.Select("*").From(mr.tableName).Where(
		goqu.C("user_id").Eq(userID),
		goqu.C("workspace_id").Eq(workspaceID),
	).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	row := executor.QueryRowContext(ctx, query)
	if err = mr.structScanRow(&membership, row); err != nil {
		return nil, err
	}
	return &membership, nil
}

func (mr *membershipRepository) Update(ctx context.Context, membership entity.Membership) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := mr.dialect.Update(mr.tableName).Set(
		membership,
	).Where(
		goqu.C("user_id").Eq(membership.UserID),
		goqu.C("workspace_id").Eq(membership.WorkspaceID),
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

func (mr *membershipRepository) SoftDelete(ctx context.Context, id string) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := mr.dialect.Update(mr.tableName).Set(
		goqu.Record{"is_deleted": true},
	).Where(
		goqu.C("id").Eq(id),
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

func (mr *membershipRepository) ListChannelMemberships(ctx context.Context, channelID string) ([]entity.Membership, error) {
	query := `
	SELECT Memberships.id, Memberships.user_id, Memberships.workspace_id, Memberships.name, Memberships.profile_image_url, Memberships.is_admin, Memberships.is_deleted
	FROM Memberships
	JOIN Membership_Channels ON Memberships.id = Membership_Channels.membership_id
	WHERE Membership_Channels.channel_id = ?;
	`

	rows, err := mr.db.QueryContext(ctx, query, channelID)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return nil, err
	}
	defer rows.Close()

	var memberships []entity.Membership
	for rows.Next() {
		var membership entity.Membership
		err = rows.Scan(
			&membership.ID,
			&membership.UserID,
			&membership.WorkspaceID,
			&membership.Name,
			&membership.ProfileImageURL,
			&membership.IsAdmin,
			&membership.IsDeleted,
		)
		if err != nil {
			log.Error("Failed to scan membership", log.Ferror(err))
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	if err = rows.Err(); err != nil {
		log.Error("Failed to iterate over rows", log.Ferror(err))
		return nil, err
	}

	return memberships, nil
}
