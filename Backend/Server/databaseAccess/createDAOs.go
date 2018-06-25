package databaseAccess

import (
	"github.com/Masterminds/squirrel"
)

func CreateUsersDao() UsersDAO {
	return &PgUsersDAO{
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func CreateLeaguesDAO() LeaguesDAO {
	return &PgLeaguesDAO{
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
