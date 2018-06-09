package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

func CreateUsersDao(db *sql.DB) UsersDAO {
	return &PgUsersDAO{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
