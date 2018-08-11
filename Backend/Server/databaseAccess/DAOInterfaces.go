package databaseAccess

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (int, string, string, error)
}

type LeaguesDAO interface {
	CreateLeague(userId int, name string, publicView, publicJoin bool) (int, error)
	IsNameInUse(name string) (bool, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	GetLeagueInformation(leagueId int) (*LeagueInformation, error)
	HasEditTeamsPermission(leagueId, userId int) (bool, error)
	GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error)
}

type TeamsDAO interface {
	CreateTeam(leagueId, userId int, name, tag string) (int, error)
	IsInfoInUse(name, tag string, leagueId int) (bool, string, error)
	GetTeamInformation(teamId, leagueId int) (*TeamInformation, error)
	DoesTeamExist(teamId, leagueId int) (bool, error)
}

type GamesDAO interface {
	CreateGame(leagueId, team1Id, team2Id, gameTime int) (int, error)
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	GetGameInformation(gameId, leagueId int) (*GameInformation, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
	ReportGame(gameId, leagueId, winnerId, scoreTeam1, scoreTeam2 int) error
}
