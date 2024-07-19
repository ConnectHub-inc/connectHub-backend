package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type channelRepository struct {
	*base[entity.Channel]
}

func NewChannelRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.ChannelRepository {
	return &channelRepository{
		base: newBase[entity.Channel](db, dialect, "Channels"),
	}
}

func (rr *channelRepository) ListMembershipChannels(ctx context.Context, membershipID string) ([]entity.Channel, error) {
	query := `
	SELECT Channels.id, Channels.workspace_id, Channels.name, Channels.description, Channels.private
	FROM Channels
	JOIN Membership_Channels ON Channels.id = Membership_Channels.channel_id
	JOIN Memberships ON Membership_Channels.membership_id = Memberships.id
	WHERE Memberships.id = ?
	`

	rows, err := rr.db.QueryContext(ctx, query, membershipID)
	if err != nil {
		log.Error("Failed to query channels", log.Ferror(err))
		return nil, err
	}
	defer rows.Close()

	var channels []entity.Channel
	for rows.Next() {
		var channel entity.Channel
		err = rows.Scan(
			&channel.ID,
			&channel.WorkspaceID,
			&channel.Name,
			&channel.Description,
			&channel.Private,
		)
		if err != nil {
			log.Error("Failed to scan channel", log.Ferror(err))
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		log.Error("Failed to iterate over rows", log.Ferror(err))
		return nil, err
	}

	return channels, nil
}
