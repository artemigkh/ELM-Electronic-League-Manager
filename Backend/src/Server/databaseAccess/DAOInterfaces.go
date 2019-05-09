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
	GetAuthenticationInformation(email string) (int, string, string, error)
	GetUserProfile(userId int) (*UserInformation, error)
	GetPermissions(leagueId, userId int) (*UserPermissions, error)
}

type LeaguesDAO interface {
	// Leagues
	GetPublicLeagueList() ([]PublicLeagueInformation, error)
	CreateLeague(userId int, name, description, game string, publicView, publicJoin bool,
		signupStart, signupEnd, leagueStart, leagueEnd int) (int, error)
	UpdateLeague(leagueId int, name, description, game string, publicView, publicJoin bool,
		signupStart, signupEnd, leagueStart, leagueEnd int) error
	GetLeagueInformation(leagueId int) (*LeagueInformation, error)
	GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error)
	GetGameSummary(leagueId int) ([]GameSummaryInformation, error)
	JoinLeague(leagueId, userId int) error
	GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error)
	SetLeaguePermissions(leagueId, userId int, administrator, createTeams, editTeams, editGames bool) error
	GetMarkdownFile(leagueId int) (string, error)
	SetMarkdownFile(leagueId int, fileName string) error
	AddRecurringAvailability(leagueId int, weekday int, timezone int,
		hour, minute, duration int, constrained bool, start, end int) (int, error)
	EditRecurringAvailability(leagueId, availabilityId int, weekday int, timezone int,
		hour, minute, duration int, constrained bool, start, end int) error
	RemoveRecurringAvailabilities(leagueId, availabilityId int) error

	// Get Information
	IsNameInUse(leagueId int, name string) (bool, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	GetLeaguePermissions(leagueId, userId int) (*LeaguePermissions, error)
	CanJoinLeague(leagueId, userId int) (bool, error)
	GetSchedulingAvailability(leagueId, availabilityId int) (*SchedulingAvailability, error)
	GetSchedulingAvailabilities(leagueId int) ([]SchedulingAvailability, error)
}

type TeamsDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, name, tag, description string) (int, error)
	CreateTeamWithIcon(leagueId, userId int, name, tag, description, small, large string) (int, error)
	DeleteTeam(leagueId, teamId int) error
	UpdateTeam(leagueId, teamId int, name, tag, description string) error
	UpdateTeamIcon(leagueId, teamId int, small, large string) error
	GetTeamInformation(leagueId, teamId int) (*TeamInformation, error)

	// Get Information
	GetTeamPermissions(teamId, userId int) (*TeamPermissions, error)
	IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error)
	DoesTeamExist(leagueId, teamId int) (bool, error)
	DoesPlayerExist(teamId, playerId int) (bool, error)
	IsTeamActive(leagueId, teamId int) (bool, error)

	// Players
	AddNewPlayer(teamId int, gameIdentifier, name, externalId, position string, mainRoster bool) (int, error)
	RemovePlayer(teamId, playerId int) error
	UpdatePlayer(teamId, playerId int, gameIdentifier, name, externalId, position string, mainRoster bool) error

	// Managers
	ChangeManagerPermissions(teamId, userId int, administrator, information, players, reportResults bool) error
}

type GamesDAO interface {
	// Games
	CreateGame(leagueId, team1Id, team2Id, gameTime int, externalId string) (int, error)
	GetGameInformation(leagueId, gameId int) (*GameInformation, error)
	ReportGame(leagueId, gameId, winnerId, scoreTeam1, scoreTeam2 int) error
	DeleteGame(leagueId, gameId int) error
	RescheduleGame(leagueId, gameId, gameTime int) error
	AddExternalId(leagueId, gameId int, externalId string) error
	GetGameInformationFromExternalId(externalId string) (*GameInformation, error)

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
