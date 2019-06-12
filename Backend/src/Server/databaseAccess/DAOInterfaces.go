package databaseAccess

import (
	"Server/lolApi"
)

/*
 * For consistency across all function signatures, the order of numerical Ids
 * should be in order of magnitude of entity:
 * first should be league, then team, then game, then user, then player
 * and all Ids should be parameters before any others
 */

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (*UserAuthenticationDTO, error)
	GetUserProfile(userId int) (*User, error)
}

type LeaguesDAO interface {
	// Modify League
	CreateLeague(userId int, leagueInfo LeagueCore) (int, error)
	UpdateLeague(leagueId int, leagueInfo LeagueCore) error
	JoinLeague(leagueId, userId int) error

	// Permissions
	SetLeaguePermissions(leagueId, userId int, permissions LeaguePermissionsCore) error
	//GetLeaguePermissions(leagueId, userId int) (*LeaguePermissionsDTO, error)
	GetTeamManagerInformation(leagueId int) ([]*TeamWithManagers, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	CanJoinLeague(leagueId, userId int) (bool, error)

	// Get Information About Leagues
	GetLeagueInformation(leagueId int) (*League, error)
	IsNameInUse(leagueId int, name string) (bool, error)
	GetPublicLeagueList() ([]*League, error)

	// Markdown
	GetMarkdownFile(leagueId int) (string, error)
	SetMarkdownFile(leagueId int, fileName string) error

	// Availabilities
	AddAvailability(leagueId int, availability AvailabilityCore) (int, error)
	GetAvailabilities(leagueId int) ([]*Availability, error)
	DeleteAvailability(leagueId, availabilityId int) error

	AddWeeklyAvailability(leagueId int, availability WeeklyAvailabilityCore) (int, error)
	GetWeeklyAvailabilities(leagueId int) ([]*WeeklyAvailability, error)
	EditWeeklyAvailability(leagueId, availabilityId int, availability WeeklyAvailabilityCore) (int, error)
	DeleteWeeklyAvailability(leagueId, availabilityId int) error
}

type TeamsDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, teamInfo TeamCore) (int, error)
	CreateTeamWithIcon(leagueId, userId int, teamInfo TeamCore, iconSmall, iconLarge string) (int, error)
	DeleteTeam(teamId int) error
	UpdateTeam(teamId int, teamInformation TeamCore) error
	UpdateTeamIcon(teamId int, small, large string) error
	GetTeamInformation(teamId int) (*TeamWithPlayers, error)
	GetAllTeamsInLeague(leagueId int) ([]*TeamWithPlayers, error)

	// Players
	CreatePlayer(leagueId, teamId int, playerInfo PlayerCore) (int, error)
	DeletePlayer(playerId int) error
	UpdatePlayer(playerId int, playerInfo PlayerCore) error

	// Get Information For Team and Player Management
	GetTeamPermissions(teamId, userId int) (*TeamPermissionsCore, error)
	IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error)
	DoesTeamExistInLeague(leagueId, teamId int) (bool, error)
	DoesPlayerExistInTeam(teamId, playerId int) (bool, error)
	IsTeamActive(leagueId, teamId int) (bool, error)

	// Managers
	ChangeManagerPermissions(teamId, userId int, teamPermissionInformation TeamPermissionsCore) error
}

type GamesDAO interface {
	// Modify Games
	CreateGame(leagueId int, gameInformation GameCreationInformation) (int, error)
	ReportGame(gameId int, gameResult GameResult) error
	DeleteGame(gameId int) error
	RescheduleGame(gameId, gameTime int) error
	AddExternalId(gameId int, externalId string) error

	// Get Game Information
	GetAllGamesInLeague(leagueId int) ([]*Game, error)
	GetGameInformation(gameId int) (*Game, error)
	GetGameInformationFromExternalId(externalId string) (*Game, error)

	// Get Information for Games Management
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
}

type LeagueOfLegendsDAO interface {
	ReportEndGameStats(leagueId, gameId, winTeamId, loseTeamId int, match *lolApi.MatchInformation) error

	GetPlayerStats(leagueId int) ([]*PlayerStats, error)
	GetTeamStats(leagueId int) ([]*TeamStats, error)
	GetChampionStats(leagueId int) ([]*ChampionStats, error)
}
