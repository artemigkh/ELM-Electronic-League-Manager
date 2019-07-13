package databaseAccess

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/snabb/isoweek"
	"time"
)

const (
	ConflictThresholdSeconds = 120
)

type PgGamesDAO struct{}

// Modify Games

func (d *PgGamesDAO) CreateGame(leagueId int, externalId *string, gameInformation GameCreationInformation) (int, error) {
	gameId := -1
	err := psql.Insert("game").
		Columns(
			"league_id",
			"external_id",
			"team1_id",
			"team2_id",
			"game_time",
			"complete",
			"winner_id",
			"loser_id",
			"score_team1",
			"score_team2",
		).
		Values(
			leagueId,
			externalId,
			gameInformation.Team1Id,
			gameInformation.Team2Id,
			gameInformation.GameTime,
			false,
			-1,
			-1,
			0,
			0,
		).
		Suffix("RETURNING \"game_id\"").
		RunWith(db).QueryRow().Scan(&gameId)

	return gameId, err
}

func (d *PgGamesDAO) ReportGame(gameId int, gameResult GameResult) error {
	fmt.Printf("%+v\n", gameResult)
	_, err := db.Exec("SELECT report_game($1,$2,$3,$4,$5)",
		gameId,
		gameResult.WinnerId,
		gameResult.LoserId,
		gameResult.ScoreTeam1,
		gameResult.ScoreTeam2,
	)
	return err
}

func (d *PgGamesDAO) ReportGameByExternalId(externalId string, gameResult GameResult) error {
	fmt.Printf("%+v\n", gameResult)
	_, err := db.Exec("SELECT report_game_by_external_id($1,$2,$3,$4,$5)",
		externalId,
		gameResult.WinnerId,
		gameResult.LoserId,
		gameResult.ScoreTeam1,
		gameResult.ScoreTeam2,
	)
	return err
}

func (d *PgGamesDAO) DeleteGame(gameId int) error {
	_, err := psql.Delete("game").
		Where("game_id = ?", gameId).
		RunWith(db).Exec()
	return err
}

func (d *PgGamesDAO) RescheduleGame(gameId, gameTime int) error {
	_, err := psql.Update("game").
		Set("game_time", gameTime).
		Where("game_id = ?", gameId).
		RunWith(db).Exec()
	return err
}

func (d *PgGamesDAO) AddExternalId(gameId int, externalId string) error {
	_, err := psql.Update("game").
		Set("external_id", externalId).
		Where("game_id = ?", gameId).
		RunWith(db).Exec()

	return err
}

// Get Game Information

func getGameInformationBuilder() squirrel.SelectBuilder {
	return psql.Select(
		"game_id",
		"external_id",
		"league_id",
		"team1_id",
		"team2_id",
		"game_time",
		"complete",
		"winner_id",
		"loser_id",
		"score_team1",
		"score_team2",
	).From("game")
}

func (d *PgGamesDAO) GetGameInformation(gameId int) (*Game, error) {
	row := getGameSelector().
		Where("game_id = ?", gameId).
		RunWith(db).QueryRow()

	return GetScannedGame(row)
}

func (d *PgGamesDAO) GetGameInformationFromExternalId(externalId string) (*Game, error) {
	row := getGameSelector().
		Where("external_id = ?", externalId).
		RunWith(db).QueryRow()

	return GetScannedGame(row)
}

func (d *PgGamesDAO) GetAllGamesInLeague(leagueId int) ([]*Game, error) {
	var games GameArray
	if err := ScanRows(getGameSelector().
		Where("game.league_id = ?", leagueId).OrderBy("game.game_time ASC"), &games); err != nil {
		return nil, err
	}

	return games.rows, nil
}

func (d *PgGamesDAO) GetSortedGames(leagueId, teamId, limit int) (*SortedGames, error) {
	var games SortedGames
	gameSelectorCompleted := getGameSelector()
	gameSelectorUpcoming := getGameSelector()

	if teamId == 0 {
		gameSelectorCompleted = gameSelectorCompleted.
			Where("game.league_id = ? AND game.complete = true", leagueId)
		gameSelectorUpcoming = gameSelectorUpcoming.
			Where("game.league_id = ? AND game.complete = false", leagueId)
	} else {
		gameSelectorCompleted = gameSelectorCompleted.
			Where("game.league_id = ? AND game.complete = true AND "+
				"(game.team1_id = ? OR game.team2_id = ?)", leagueId, teamId, teamId)
		gameSelectorUpcoming = gameSelectorUpcoming.
			Where("game.league_id = ? AND game.complete = false AND "+
				"(game.team1_id = ? OR game.team2_id = ?)", leagueId, teamId, teamId)
	}
	gameSelectorCompleted = gameSelectorCompleted.OrderBy("game.game_time ASC")
	gameSelectorUpcoming = gameSelectorUpcoming.OrderBy("game.game_time ASC")

	if limit > 0 {
		gameSelectorCompleted = gameSelectorCompleted.Limit(uint64(limit))
		gameSelectorUpcoming = gameSelectorUpcoming.Limit(uint64(limit))
	}

	var completedGames GameArray
	if err := ScanRows(gameSelectorCompleted, &completedGames); err != nil {
		return nil, err
	}

	var upcomingGames GameArray
	if err := ScanRows(gameSelectorUpcoming, &upcomingGames); err != nil {
		return nil, err
	}

	games.CompletedGames = completedGames.rows
	games.UpcomingGames = upcomingGames.rows

	return &games, nil
}

// Get Information for Games Management

func (d *PgGamesDAO) DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error) {
	//check if any game of each team is within the threshold of another scheduled game
	//team1Games, err := getGamesOfTeam(team1Id)
	//if err != nil {
	//	return false, err
	//}
	//
	//for _, game := range team1Games {
	//	if math.Abs(float64(gameTime)-float64(game.GameTime)) < ConflictThresholdSeconds {
	//		return true, nil
	//	}
	//}
	//
	//team2Games, err := getGamesOfTeam(team2Id)
	//if err != nil {
	//	return false, err
	//}
	//
	//for _, game := range team2Games {
	//	if math.Abs(float64(gameTime)-float64(game.GameTime)) < ConflictThresholdSeconds {
	//		return true, nil
	//	}
	//}

	return false, nil
}

func getTeamsInGame(gameId int) (int, int, error) {
	team1 := -1
	team2 := -1

	err := psql.Select("team1_id", "team2_id").
		From("game").
		Where("game_id = ?", gameId).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	return team1, team2, err
}

func (d *PgGamesDAO) HasReportResultPermissions(leagueId, gameId, userId int) (bool, error) {
	//check if user has league editResults permission
	canReport := false
	err := psql.Select("edit_games").
		From("league_permissions").
		Where("league_id = ? AND user_id = ?", leagueId, userId).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	//check if user has team reportResult permissions on one of the two teams
	team1Id, team2Id, err := getTeamsInGame(gameId)

	err = psql.Select("games").
		From("team_permissions").
		Where("user_id = ? AND (team_id = ? OR team_id = ?)", userId, team1Id, team2Id).
		RunWith(db).QueryRow().Scan(&canReport)

	return canReport, err
}

type CompetitionWeek struct {
	WeekStart int     `json:"weekStart"`
	Games     []*Game `json:"games"`
}

func (d *PgGamesDAO) GetGamesByWeek(leagueId, timeZone int) ([]*CompetitionWeek, error) {
	games, err := d.GetAllGamesInLeague(leagueId)
	if err != nil {
		return nil, err
	}

	competitionWeeks := make([]*CompetitionWeek, 0)
	if len(games) == 0 {
		return competitionWeeks, nil
	}

	// Create initial objects for the first week of games
	year, week := time.Unix(int64(games[0].GameTime), 0).ISOWeek()
	weekStart := isoweek.StartTime(year, week, time.FixedZone("", timeZone))
	competitionWeek := &CompetitionWeek{
		WeekStart: int(weekStart.Unix()),
		Games:     make([]*Game, 0),
	}
	competitionWeeks = append(competitionWeeks, competitionWeek)

	// Add all games to a week struct, creating a new week as necessary
	for _, game := range games {
		// Make new week if game is after the end of current week
		if time.Unix(int64(game.GameTime), 0).After(weekStart.AddDate(0, 0, 7)) {
			weekStart = weekStart.AddDate(0, 0, 7)
			competitionWeek = &CompetitionWeek{
				WeekStart: int(weekStart.Unix()),
				Games:     make([]*Game, 0),
			}
			competitionWeeks = append(competitionWeeks, competitionWeek)
		}
		competitionWeek.Games = append(competitionWeek.Games, game)
	}

	return competitionWeeks, nil
}
