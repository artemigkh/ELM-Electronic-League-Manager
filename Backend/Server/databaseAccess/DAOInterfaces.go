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
	CreateLeague(userId int, name string, publicView, publicJoin bool) (int, error)
	IsNameInUse(name string) (bool, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	GetLeagueInformation(leagueId int) (*LeagueInformation, error)
	HasEditTeamsPermission(leagueId, userId int) (bool, error)
	HasCreateTeamsPermission(leagueId, userId int) (bool, error)
	GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error)
	CanJoinLeague(leagueId, userId int) (bool, error)
	JoinLeague(leagueId, userId int) error
	IsLeagueAdmin(leagueId, userId int) (bool, error)
	GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error)
}

type TeamsDAO interface {
	CreateTeam(leagueId, userId int, name, tag string) (int, error)
	IsInfoInUse(leagueId int, name, tag string) (bool, string, error)
	GetTeamInformation(leagueId, teamId int) (*TeamInformation, error)
	DoesTeamExist(leagueId, teamId int) (bool, error)
	HasPlayerEditPermissions(leagueId, teamId, userId int) (bool, error)
	AddNewPlayer(teamId int, gameIdentifier, name string, mainRoster bool) (int, error)
}

type GamesDAO interface {
	CreateGame(leagueId, team1Id, team2Id, gameTime int) (int, error)
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	GetGameInformation(leagueId, gameId int) (*GameInformation, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
	ReportGame(leagueId, gameId, winnerId, scoreTeam1, scoreTeam2 int) error
}
