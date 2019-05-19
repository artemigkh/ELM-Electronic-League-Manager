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

//type ElmDAO interface {
//	CreateLeague(userId int, leagueInfo LeagueDTO) (int, error)
//
//	GetLeagueInformation(leagueId int) (*LeagueDTO, error)
//	GetPublicLeagueList() ([]*LeagueDTO, error)
//	GetTeamSummary(leagueId int) ([]*TeamDTO, error)
//	GetGameSummary(leagueId int) ([]*GameSummaryInformationDTO, error)
//	GetTeamManagerInformation(leagueId int) ([]TeamManagerDTO, error)
//
//	UpdateLeague(leagueInfo LeagueDTO) error
//
//	JoinLeague(leagueId, userId int) error
//}

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (*UserAuthenticationDTO, error)
	GetUserProfile(userId int) (*UserDTO, error)
	GetPermissions(leagueId, userId int) (*UserPermissionsDTO, error)
}

type LeaguesDAO interface {
	// Modify League
	CreateLeague(userId int, leagueInfo LeagueDTO) (int, error)
	UpdateLeague(leagueInfo LeagueDTO) error
	JoinLeague(leagueId, userId int) error

	// Permissions
	SetLeaguePermissions(leagueId, userId int, permissions LeaguePermissionsDTO) error
	GetLeaguePermissions(leagueId, userId int) (*LeaguePermissionsDTO, error)
	GetTeamManagerInformation(leagueId int) ([]*TeamManagerDTO, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	CanJoinLeague(leagueId, userId int) (bool, error)

	// Get Information About Leagues
	GetLeagueInformation(leagueId int) (*LeagueDTO, error)
	IsNameInUse(leagueId int, name string) (bool, error)
	GetPublicLeagueList() ([]*LeagueDTO, error)

	// Get Information About Entities in a League
	GetTeamSummary(leagueId int) ([]*TeamDTO, error)
	GetGameSummary(leagueId int) ([]*GameDTO, error)

	// Markdown
	GetMarkdownFile(leagueId int) (string, error)
	SetMarkdownFile(leagueId int, fileName string) error

	// Availabilities
	AddRecurringAvailability(leagueId int, availability SchedulingAvailabilityDTO) (int, error)
	EditRecurringAvailability(leagueId int, availability SchedulingAvailabilityDTO) error
	RemoveRecurringAvailabilities(leagueId, availabilityId int) error
	GetSchedulingAvailability(leagueId, availabilityId int) (*SchedulingAvailabilityDTO, error)
	GetSchedulingAvailabilities(leagueId int) ([]*SchedulingAvailabilityDTO, error)
}

type TeamsDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, teamInfo TeamDTO) (int, error)
	CreateTeamWithIcon(leagueId, userId int, teamInfo TeamDTO) (int, error)
	DeleteTeam(leagueId, teamId int) error
	UpdateTeam(leagueId, teamInformation TeamDTO) error
	UpdateTeamIcon(leagueId, teamId int, small, large string) error
	GetTeamInformation(leagueId, teamId int) (*TeamDTO, error)

	// Players
	AddNewPlayer(playerInfo PlayerDTO) (int, error)
	RemovePlayer(teamId, playerId int) error
	UpdatePlayer(playerInfo PlayerDTO) error

	// Get Information For Team and Player Management
	GetTeamPermissions(teamId, userId int) (*TeamPermissionsDTO, error)
	IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error)
	DoesTeamExist(leagueId, teamId int) (bool, error)
	DoesPlayerExist(teamId, playerId int) (bool, error)
	IsTeamActive(leagueId, teamId int) (bool, error)

	// Managers
	ChangeManagerPermissions(teamId, userId int, teamPermissionInformation TeamPermissionsDTO) error
}

type GamesDAO interface {
	// Modify Games
	CreateGame(gameInformation GameDTO) (int, error)
	ReportGame(gameInfo GameDTO) error
	DeleteGame(leagueId, gameId int) error
	RescheduleGame(leagueId, gameId, gameTime int) error
	AddExternalId(leagueId, gameId int, externalId string) error

	// Get Game Information
	GetGameInformation(leagueId, gameId int) (*GameDTO, error)
	GetGameInformationFromExternalId(externalId string) (*GameDTO, error)

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
