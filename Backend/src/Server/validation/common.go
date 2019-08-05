package validation

import "Server/dataModel"

type AccessType int

const (
	View   AccessType = iota
	Edit   AccessType = iota
	Create AccessType = iota
	Delete AccessType = iota
)

// checks availability and permissions of user regarding data type
type Access interface {
	League(accessType AccessType, permissions *dataModel.UserWithPermissions,
		leagueDao dataModel.LeagueDAO, leagueId int) (bool, error)
	Team(accessType AccessType, permissions *dataModel.UserWithPermissions,
		teamDao dataModel.TeamDAO, leagueDao dataModel.LeagueDAO,
		leagueId, teamId int) (bool, error)
	Player(accessType AccessType, permissions *dataModel.UserWithPermissions,
		teamDao dataModel.TeamDAO, leagueDao dataModel.LeagueDAO,
		leagueId, teamId, playerId int) (bool, error)
	Game(accessType AccessType, permissions *dataModel.UserWithPermissions,
		gameDao dataModel.GameDAO, leagueId, gameId int) (bool, error)
	Availability(accessType AccessType, permissions *dataModel.UserWithPermissions,
		leagueDao dataModel.LeagueDAO, leagueId, availabilityId int) (bool, error)
}
type AccessChecker struct{}
