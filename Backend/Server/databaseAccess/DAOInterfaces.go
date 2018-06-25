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
	GetLeagueInformation(userID int) (*LeagueInformation, error)
}