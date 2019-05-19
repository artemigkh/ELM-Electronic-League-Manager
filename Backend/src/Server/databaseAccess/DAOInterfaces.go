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
//	CreateLeague(userId int, leagueInfo LeagueInformationDTO) (int, error)
//
//	GetLeagueInformation(leagueId int) (*LeagueInformationDTO, error)
//	GetPublicLeagueList() ([]*LeagueInformationDTO, error)
//	GetTeamSummary(leagueId int) ([]*TeamDTO, error)
//	GetGameSummary(leagueId int) ([]*GameSummaryInformationDTO, error)
//	GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error)
//
//	UpdateLeague(leagueInfo LeagueInformationDTO) error
//
//	JoinLeague(leagueId, userId int) error
//}

type UsersDAO interface {
	CreateUser(email, salt, hash string) error
	IsEmailInUse(email string) (bool, error)
	GetAuthenticationInformation(email string) (int, string, string, error)
	GetUserProfile(userId int) (*UserInformation, error)
	GetPermissions(leagueId, userId int) (*UserPermissions, error)
}

type LeaguesDAO interface {
	// Leagues
	GetPublicLeagueList() ([]*LeagueInformationDTO, error)
	CreateLeague(userId int, leagueInfo LeagueInformationDTO) (int, error)
	UpdateLeague(leagueInfo LeagueInformationDTO) error
	GetLeagueInformation(leagueId int) (*LeagueInformationDTO, error)
	GetTeamSummary(leagueId int) ([]*TeamDTO, error)
	GetGameSummary(leagueId int) ([]*GameDTO, error)
	JoinLeague(leagueId, userId int) error
	GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error)
	SetLeaguePermissions(leagueId, userId int, permissions LeaguePermissionsDTO) error
	GetMarkdownFile(leagueId int) (string, error)
	SetMarkdownFile(leagueId int, fileName string) error

	// Availabilities
	AddRecurringAvailability(leagueId int, availability SchedulingAvailabilityDTO) (int, error)
	EditRecurringAvailability(leagueId int, availability SchedulingAvailabilityDTO) error
	RemoveRecurringAvailabilities(leagueId, availabilityId int) error
	GetSchedulingAvailability(leagueId, availabilityId int) (*SchedulingAvailabilityDTO, error)
	GetSchedulingAvailabilities(leagueId int) ([]*SchedulingAvailabilityDTO, error)

	// Get Information
	IsNameInUse(leagueId int, name string) (bool, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	GetLeaguePermissions(leagueId, userId int) (*LeaguePermissionsDTO, error)
	CanJoinLeague(leagueId, userId int) (bool, error)
}

type TeamsDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, teamInfo TeamDTO) (int, error)
	CreateTeamWithIcon(leagueId, userId int, teamInfo TeamDTO) (int, error)
	DeleteTeam(leagueId, teamId int) error
	UpdateTeam(leagueId, teamInformation TeamDTO) error
	UpdateTeamIcon(leagueId, teamId int, small, large string) error
	GetTeamInformation(leagueId, teamId int) (*TeamDTO, error)

	// Get Information
	GetTeamPermissions(teamId, userId int) (*TeamPermissionsDTO, error)
	IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error)
	DoesTeamExist(leagueId, teamId int) (bool, error)
	DoesPlayerExist(teamId, playerId int) (bool, error)
	IsTeamActive(leagueId, teamId int) (bool, error)

	// Players
	AddNewPlayer(playerInfo PlayerDTO) (int, error)
	RemovePlayer(teamId, playerId int) error
	UpdatePlayer(playerInfo PlayerDTO) error

	// Managers
	ChangeManagerPermissions(teamId, userId int, teamPermissionInformation TeamPermissionsDTO) error
}

type GamesDAO interface {
	// Games
	CreateGame(leagueId, team1Id, team2Id, gameTime int, externalId string) (int, error)
	GetGameInformation(leagueId, gameId int) (*GameDTO, error)
	ReportGame(leagueId, gameId, winnerId, scoreTeam1, scoreTeam2 int) error
	DeleteGame(leagueId, gameId int) error
	RescheduleGame(leagueId, gameId, gameTime int) error
	AddExternalId(leagueId, gameId int, externalId string) error
	GetGameInformationFromExternalId(externalId string) (*GameDTO, error)

	// Get Information
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
}

type LeagueOfLegendsDAO interface {
	ReportEndGameStats(leagueId, gameId, winTeamId, loseTeamId int, match *lolApi.MatchInformation) error

	GetPlayerStats(leagueId int) ([]*PlayerStats, error)
	GetTeamStats(leagueId int) ([]*TeamStats, error)
	GetChampionStats(leagueId int) ([]*ChampionStats, error)
}
