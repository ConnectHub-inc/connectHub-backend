package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
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

func (rr *roomRepository) ListMembershipRooms(ctx context.Context, membershipID string) ([]entity.Room, error) {
	query := `
	SELECT Rooms.id, Rooms.workspace_id, Rooms.name, Rooms.description, Rooms.private
	FROM Rooms
	JOIN Membership_Rooms ON Rooms.id = Membership_Rooms.room_id
	JOIN Memberships ON Membership_Rooms.membership_id = Memberships.id
	WHERE Memberships.id = ?
	`

	rows, err := rr.db.QueryContext(ctx, query, membershipID)
	if err != nil {
		log.Error("Failed to query rooms", log.Ferror(err))
		return nil, err
	}
	defer rows.Close()

	var rooms []entity.Room
	for rows.Next() {
		var room entity.Room
		err = rows.Scan(
			&room.ID,
			&room.WorkspaceID,
			&room.Name,
			&room.Description,
			&room.Private,
		)
		if err != nil {
			log.Error("Failed to scan room", log.Ferror(err))
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		log.Error("Failed to iterate over rows", log.Ferror(err))
		return nil, err
	}

	return rooms, nil
}
