package databaseAccess

func CreateUsersDao() UsersDAO {
	return &PgUsersDAO{}
}

func CreateLeaguesDAO() LeaguesDAO {
	return &PgLeaguesDAO{}
}

func CreateTeamsDAO() TeamsDAO {
	return &PgTeamsDAO{}
}

func CreateGamesDAO() GamesDAO {
	return &PgGamesDAO{}
}
