package databaseAccess

import (
	"github.com/Masterminds/squirrel"
)

const (
	ConflictThresholdSeconds = 120
)

type PgGamesDAO struct{}

// Modify Games

func (d *PgGamesDAO) CreateGame(leagueId int, gameInformation GameCreationInformation) (int, error) {
	gameId := -1
	err := psql.Insert("game").
		Columns(
			"league_id",
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
	_, err := db.Exec("SELECT report_game($1,$2,$3,$4,$5)",
		gameId,
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
		Where("game.league_id = ?", leagueId), &games); err != nil {
		return nil, err
	}

	return games.rows, nil
}

// Get Information for Games Management

func getGamesOfTeam(teamId int) ([]*Game, error) {
	//var games GameDTOArray
	//if err := ScanRows(psql.Select(
	//	"game_id",
	//	"external_id",
	//	"league_id",
	//	"team1_id",
	//	"team2_id",
	//	"game_time",
	//	"complete",
	//	"winner_id",
	//	"loser_id",
	//	"score_team1",
	//	"score_team2",
	//).
	//	From("game").
	//	Where("team1_id = ? OR team2_id = ?", teamId, teamId), &games); err != nil {
	//	return nil, err
	//}
	//
	//return games.rows, nil
	return nil, nil
}

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
