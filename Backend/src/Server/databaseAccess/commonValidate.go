package databaseAccess

type AccessType int

const (
	View   AccessType = iota
	Edit   AccessType = iota
	Create AccessType = iota
	Delete AccessType = iota
)

// checks availability and permissions of user regarding data type
type Access interface {
	League(accessType AccessType, leagueId, userId int) (bool, error)
	Team(accessType AccessType, leagueId, teamId, userId int) (bool, error)
	Player(accessType AccessType, leagueId, teamId, playerId, userId int) (bool, error)
	Game(accessType AccessType, leagueId, gameId, userId int) (bool, error)
	Report(leagueId, gameId, userId int) (bool, error)
}
type AccessChecker struct{}
