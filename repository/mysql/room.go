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
	var membership entity.Membership
	userID, workspaceID, err := membership.SplitMembershipID(membershipID)
	if err != nil {
		log.Error("Failed to split membership ID", log.Ferror(err))
		return nil, err
	}

	query := `
	SELECT Rooms.id, Rooms.workspace_id, Rooms.name, Rooms.description, Rooms.private
	FROM Rooms
	JOIN Workspaces ON Rooms.workspace_id = Workspaces.id
	JOIN Memberships ON Workspaces.id = Memberships.workspace_id
	JOIN Membership_Rooms ON Rooms.id = Membership_Rooms.room_id
	WHERE Memberships.user_id = ?
	  AND Memberships.workspace_id = ?
  	  AND Membership_Rooms.user_id = ?;
	`

	rows, err := rr.db.QueryContext(ctx, query, userID, workspaceID, userID)
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
