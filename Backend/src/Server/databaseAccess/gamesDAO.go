package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"math"
)

const (
	ConflictThresholdSeconds = 120
)

type PgGamesDAO struct{}

// Modify Games

func (d *PgGamesDAO) CreateGame(gameInformation GameDTO) (int, error) {
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
			gameInformation.LeagueId,
			gameInformation.ExternalId,
			gameInformation.Team1Id,
			gameInformation.Team2Id,
			gameInformation.GameTime,
			false,
			-1,
			-1,
			0,
			0,
		).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&gameId)

	return gameId, err
}

func (d *PgGamesDAO) ReportGame(gameInfo GameDTO) error {
	_, err := db.Exec("SELECT report_game($1,$2,$3,$4,$5)",
		gameInfo.Id,
		gameInfo.WinnerId,
		gameInfo.LoserId,
		gameInfo.ScoreTeam1,
		gameInfo.ScoreTeam2,
	)
	return err
}

func (d *PgGamesDAO) DeleteGame(leagueId, gameId int) error {
	_, err := psql.Delete("game").
		Where("id = ? AND league_id = ?", gameId, leagueId).
		RunWith(db).Exec()
	return err
}

func (d *PgGamesDAO) RescheduleGame(leagueId, gameId, gameTime int) error {
	_, err := psql.Update("game").
		Set("game_time", gameTime).
		Where("id = ? AND league_id = ?", gameId, leagueId).RunWith(db).Exec()
	return err
}

func (d *PgGamesDAO) AddExternalId(leagueId, gameId int, externalId string) error {
	_, err := psql.Update("game").
		Set("external_id", externalId).
		Where("id = ? AND league_id", gameId, leagueId).
		RunWith(db).Exec()

	return err
}

// Get Game Information

func getGameInformationBuilder() squirrel.SelectBuilder {
	return psql.Select(
		"id",
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

func (d *PgGamesDAO) GetGameInformation(leagueId, gameId int) (*GameDTO, error) {
	return GetScannedGameDTO(getGameInformationBuilder().
		Where("id = ? AND league_id = ?", gameId, leagueId).
		RunWith(db).QueryRow())
}

func (d *PgGamesDAO) GetGameInformationFromExternalId(externalId string) (*GameDTO, error) {
	return GetScannedGameDTO(getGameInformationBuilder().
		Where("external_id = ?", externalId).
		RunWith(db).QueryRow())
}

// Get Information for Games Management

func (d *PgGamesDAO) DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error) {
	//check if any game of each team is within the threshold of another scheduled game
	team1Games, err := getGamesOfTeam(team1Id)
	if err != nil {
		return false, err
	}

	for _, game := range team1Games {
		if math.Abs(float64(gameTime)-float64(game.GameTime)) < ConflictThresholdSeconds {
			return true, nil
		}
	}

	team2Games, err := getGamesOfTeam(team2Id)
	if err != nil {
		return false, err
	}

	for _, game := range team2Games {
		if math.Abs(float64(gameTime)-float64(game.GameTime)) < ConflictThresholdSeconds {
			return true, nil
		}
	}

	return false, nil
}

func getTeamsInGame(gameId, leagueId int) (int, int, error) {
	team1 := -1
	team2 := -1

	err := psql.Select("team1_id", "team2_id").
		From("game").
		Where("id = ? AND league_id = ?", gameId, leagueId).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	return team1, team2, err
}

func (d *PgGamesDAO) HasReportResultPermissions(leagueId, gameId, userId int) (bool, error) {
	//check if user has league editResults permission
	canReport := false
	err := psql.Select("edit_games").
		From("league_permissions").
		Where("user_id = ? AND league_id = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	//check if user has team reportResult permissions on one of the two teams
	team1Id, team2Id, err := getTeamsInGame(gameId, leagueId)

	err = psql.Select("report_results").
		From("team_permissions").
		Where("user_id = ? AND (team_id = ? OR team_id = ?)", userId, team1Id, team2Id).
		RunWith(db).QueryRow().Scan(&canReport)

	return canReport, err
}

func getGamesOfTeam(teamId int) ([]*GameDTO, error) {
	var games GameDTOArray
	if err := ScanRows(psql.Select(
		"id",
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
	).
		From("game").
		Where("team1_id = ? OR team2_id = ?", teamId, teamId), &games); err != nil {
		return nil, err
	}

	return games.rows, nil
}
