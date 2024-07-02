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

func (rr *roomRepository) ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error) {
	query := `
	SELECT Rooms.id, Rooms.workspace_id, Rooms.name, Rooms.description, Rooms.private
	FROM Rooms
	JOIN Workspaces ON Rooms.workspace_id = Workspaces.id
	JOIN User_Workspaces ON Workspaces.id = User_Workspaces.workspace_id
	JOIN User_Rooms ON Rooms.id = User_Rooms.room_id
	WHERE User_Workspaces.user_id = ?
	  AND User_Workspaces.workspace_id = ?
  	  AND User_Rooms.user_workspace_id = CONCAT(User_Workspaces.user_id, '_', User_Workspaces.workspace_id);
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
