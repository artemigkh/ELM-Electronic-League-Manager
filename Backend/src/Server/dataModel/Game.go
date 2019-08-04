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

func (gameTime *GameTime) Validate(leagueId, gameId int, leagueDao LeagueDAO, teamDao TeamDAO, gameDao GameDAO) (bool, string, error) {
	gameInformation, err := gameDao.GetGameInformation(gameId)
	if err != nil {
		return false, "", err
	} else {
		return GameCreationInformation{
			Team1Id:  gameInformation.Team1.TeamId,
			Team2Id:  gameInformation.Team2.TeamId,
			GameTime: gameTime.GameTime,
		}.ValidateReschedule(leagueId, gameId, leagueDao, teamDao, gameDao)
	}
}

type GameCreationInformation struct {
	Team1Id  int `json:"team1Id"`
	Team2Id  int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

func (game GameCreationInformation) validate(leagueId, gameId int, leagueDao LeagueDAO, teamDao TeamDAO, gameDao GameDAO) (bool, string, error) {
	return validate(
		game.differentTeams(),
		game.teamsExist(leagueId, teamDao),
		game.noConflict(leagueId, gameId, gameDao),
		validateDuringLeague(leagueId, game.GameTime, leagueDao, GameNotDuringLeague))
}

func (game GameCreationInformation) Validate(leagueId int, leagueDao LeagueDAO, teamDao TeamDAO, gameDao GameDAO) (bool, string, error) {
	return game.validate(leagueId, 0, leagueDao, teamDao, gameDao)
}

func (game GameCreationInformation) ValidateReschedule(leagueId, gameId int, leagueDao LeagueDAO, teamDao TeamDAO, gameDao GameDAO) (bool, string, error) {
	return game.validate(leagueId, gameId, leagueDao, teamDao, gameDao)
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

func (game *GameCreationInformation) noConflict(leagueId, gameId int, gameDao GameDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		allGames, err := gameDao.GetAllGamesInLeague(leagueId)
		if err != nil {
			*errorDest = err
			return false
		}
		for _, ExistingGame := range allGames {
			if gameId != ExistingGame.GameId &&
				game.GameTime == ExistingGame.GameTime &&
				(game.Team1Id == ExistingGame.Team1.TeamId ||
					game.Team1Id == ExistingGame.Team2.TeamId ||
					game.Team2Id == ExistingGame.Team1.TeamId ||
					game.Team2Id == ExistingGame.Team2.TeamId) {
				*problemDest = GameConflict
				return false
			}
		}
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

func (gameResult *GameResult) Validate(gameId int, gameDao GameDAO) (bool, string, error) {
	gameInformation, err := gameDao.GetGameInformation(gameId)
	if err != nil {
		return false, "", err
	} else {
		return validate(gameResult.hasReportedTeams(gameInformation))
	}
}

func (gameResult *GameResult) ValidateByExternalId(externalGameId string, gameDao GameDAO) (bool, string, error) {
	gameInformation, err := gameDao.GetGameInformationFromExternalId(externalGameId)
	if err != nil {
		return false, "", err
	} else {
		return validate(gameResult.hasReportedTeams(gameInformation))
	}
}

func (gameResult *GameResult) hasReportedTeams(gameInformation *Game) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if (gameInformation.Team1.TeamId == gameResult.WinnerId &&
			gameInformation.Team2.TeamId == gameResult.LoserId) ||
			(gameInformation.Team2.TeamId == gameResult.WinnerId &&
				gameInformation.Team1.TeamId == gameResult.LoserId) {
			return true
		} else {
			*problemDest = GameDoesNotContainTeams
			return false
		}
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
