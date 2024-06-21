package mysql

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type userRepository struct {
	*base[entity.User]
}

func NewUserRepository(db *sql.DB, dialect *goqu.DialectWrapper) repository.UserRepository {
	return &userRepository{
		base: newBase[entity.User](db, dialect, "Users"),
	}
}

func (ur *userRepository) ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.User, error) {
	query := `
	SELECT Users.id, Users.name
	FROM Users
	JOIN User_Workspaces ON Users.id = User_Workspaces.user_id
	WHERE User_Workspaces.workspace_id = ?;
	`

	rows, err := ur.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		log.Error("Failed to query workspace users", log.Ferror(err))
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			log.Error("Failed to scan user", log.Ferror(err))
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
