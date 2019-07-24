package dataModel

type GameDAO interface {
	// Modify Games
	CreateGame(leagueId int, externalId *string, gameInformation GameCreationInformation) (int, error)
	ReportGame(gameId int, gameResult GameResult) error
	ReportGameByExternalId(externalId string, gameResult GameResult) (int, int, error)
	DeleteGame(gameId int) error
	RescheduleGame(gameId, gameTime int) error
	AddExternalId(gameId int, externalId string) error

	// Get Game Information
	GetAllGamesInLeague(leagueId int) ([]*Game, error)
	GetSortedGames(leagueId, teamId, limit int) (*SortedGames, error)
	GetGamesByWeek(leagueId, timeZone int) ([]*CompetitionWeek, error)
	GetGameInformation(gameId int) (*Game, error)
	GetGameInformationFromExternalId(externalId string) (*Game, error)
	DoesGameExistInLeague(leagueId, gameId int) (bool, error)

	// Get Information for Games Management
	DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error)
	HasReportResultPermissions(leagueId, gameId, userId int) (bool, error)
}

type GameTime struct {
	GameTime int `json:"gameTime"`
}

func (gameTime *GameTime) Validate(leagueId, gameId int) (bool, string, error) {
	return validate(
		gameTime.noConflict(),
		gameTime.duringLeague())
}

func (gameTime *GameTime) noConflict() ValidateFunc {
	return func(_ *string, _ *error) bool {
		// TODO
		return true
	}
}

func (gameTime *GameTime) duringLeague() ValidateFunc {
	return func(_ *string, _ *error) bool {
		// TODO
		return true
	}
}

type GameCreationInformation struct {
	Team1Id  int `json:"team1Id"`
	Team2Id  int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

func (game *GameCreationInformation) Validate(leagueId int, teamDao TeamDAO) (bool, string, error) {
	return validate(
		game.differentTeams(),
		game.teamsExist(leagueId, teamDao),
		game.noConflict(),
		game.duringLeague())
}

func (game *GameCreationInformation) differentTeams() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if game.Team1Id == game.Team2Id {
			*problemDest = TeamsMustBeDifferent
			return false
		} else {
			return true
		}
	}
}

func (game *GameCreationInformation) teamsExist(leagueId int, teamDao TeamDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		team1Exists, team1QueryErr := teamDao.DoesTeamExistInLeague(leagueId, game.Team1Id)
		team2Exists, team2QueryErr := teamDao.DoesTeamExistInLeague(leagueId, game.Team2Id)
		if team1QueryErr != nil {
			*errorDest = team1QueryErr
		} else if team2QueryErr != nil {
			*errorDest = team2QueryErr
		} else if !(team1Exists && team2Exists) {
			*problemDest = TeamInGameDoesNotExist
		} else {
			valid = true
		}
		return valid
	}
}

func (game *GameCreationInformation) noConflict() ValidateFunc {
	return func(_ *string, _ *error) bool {
		// TODO
		return true
	}
}

func (game *GameCreationInformation) duringLeague() ValidateFunc {
	return func(_ *string, _ *error) bool {
		// TODO
		return true
	}
}

type GameCore struct {
	GameTime int         `json:"gameTime"`
	Team1    TeamDisplay `json:"team1"`
	Team2    TeamDisplay `json:"team2"`
}

type GameResult struct {
	WinnerId   int `json:"winnerId"`
	LoserId    int `json:"loserId"`
	ScoreTeam1 int `json:"scoreTeam1"`
	ScoreTeam2 int `json:"scoreTeam2"`
}

func (gameResult *GameResult) Validate(leagueId int) (bool, string, error) {
	return validate(
		gameResult.hasReportedTeams())
}

func (gameResult *GameResult) hasReportedTeams() ValidateFunc {
	return func(_ *string, _ *error) bool {
		// TODO
		return true
	}
}

type Game struct {
	GameId     int         `json:"gameId"`
	GameTime   int         `json:"gameTime"`
	Team1      TeamDisplay `json:"team1"`
	Team2      TeamDisplay `json:"team2"`
	WinnerId   int         `json:"winnerId"`
	LoserId    int         `json:"loserId"`
	ScoreTeam1 int         `json:"scoreTeam1"`
	ScoreTeam2 int         `json:"scoreTeam2"`
	Complete   bool        `json:"complete"`
}

type SortedGames struct {
	CompletedGames []*Game `json:"completedGames"`
	UpcomingGames  []*Game `json:"upcomingGames"`
}

type CompetitionWeek struct {
	WeekStart int     `json:"weekStart"`
	Games     []*Game `json:"games"`
}
