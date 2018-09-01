package databaseAccess

/*
 * For consistency across all function signatures, the order of numerical Ids
 * should be in order of magnitude of entity:
 * first should be league, then team, then game, then user, then player
 * and all Ids should be parameters before any others
 */

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (int, string, string, error)
	GetUserProfile(userId int) (*UserInformation, error)
}

type LeaguesDAO interface {
	// Leagues
	CreateLeague(userId int, name string, publicView, publicJoin bool) (int, error)
	GetLeagueInformation(leagueId int) (*LeagueInformation, error)
	GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error)
	JoinLeague(leagueId, userId int) error
	GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error)

	// Get Information
	IsNameInUse(name string) (bool, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	HasEditTeamsPermission(leagueId, userId int) (bool, error)
	HasEditSchedulePermission(leagueId, userId int) (bool, error)
	HasCreateTeamsPermission(leagueId, userId int) (bool, error)
	CanJoinLeague(leagueId, userId int) (bool, error)
	IsLeagueAdmin(leagueId, userId int) (bool, error)
}

type TeamsDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, name, tag string) (int, error)
	GetTeamInformation(leagueId, teamId int) (*TeamInformation, error)

	// Get Information
	IsInfoInUse(leagueId int, name, tag string) (bool, string, error)
	HasPlayerEditPermissions(leagueId, teamId, userId int) (bool, error)
	DoesTeamExist(leagueId, teamId int) (bool, error)
	DoesPlayerExist(teamId, playerId int) (bool, error)

	// Players
	AddNewPlayer(teamId int, gameIdentifier, name string, mainRoster bool) (int, error)
	RemovePlayer(teamId, playerId int) error
}

type GamesDAO interface {
	// Games
	CreateGame(leagueId, team1Id, team2Id, gameTime int) (int, error)
	GetGameInformation(leagueId, gameId int) (*GameInformation, error)
	ReportGame(leagueId, gameId, winnerId, scoreTeam1, scoreTeam2 int) error
	DeleteGame(leagueId, gameId int) error

	// Get Information
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
}
