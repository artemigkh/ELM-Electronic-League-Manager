package databaseAccess

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (int, string, string, error)
}

type LeaguesDAO interface {
	CreateLeague(userID int, name string, publicView, publicJoin bool) (int, error)
	IsNameInUse(name string) (bool, error)
	IsLeagueViewable(leagueID, userID int) (bool, error)
	GetLeagueInformation(leagueID int) (*LeagueInformation, error)
	HasEditTeamsPermission(leagueID, userID int) (bool, error)
}

type TeamsDAO interface {
	CreateTeam(leagueID, userID int, name, tag string) (int, error)
	IsInfoInUse(name, tag string, leagueID int) (bool, string, error)
}