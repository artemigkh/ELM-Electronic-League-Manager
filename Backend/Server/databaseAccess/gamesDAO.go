package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"math"
	"database/sql"
	"errors"
)

const (
	SCHEDULING_CONFLICT_THRESHOLD_SECONDS = 120
)

type GameInformation struct {
	Id int `json:"id"`
	LeagueID int `json:"leagueId"`
	Team1ID int `json:"team1Id"`
	Team2ID int `json:"team2Id"`
	GameTime int `json:"gameTime"`
	Complete bool `json:"complete"`
	WinnerID int `json:"winnerId"`
	ScoreTeam1 int `json:"scoreTeam1"`
	ScoreTeam2 int `json:"scoreTeam2"`
}

type PgGamesDAO struct {}

func getGamesOfTeam(teamId int) ([]GameInformation, error) {
	rows, err := psql.Select("*").From("games").
		Where(squirrel.Or{squirrel.Eq{"team1ID": teamId}, squirrel.Eq{"team2ID": teamId}}).
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameInformation
	var game GameInformation

	for rows.Next() {
		err := rows.Scan(&game.Id, &game.LeagueID, &game.Team1ID, &game.Team2ID, &game.GameTime,
			&game.Complete, &game.WinnerID, &game.ScoreTeam1, &game.ScoreTeam2)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return games, nil
}

func getTeamsInGame(gameID, leagueID int) (int, int, error) {
	var (
		team1 int
		team2 int
	)

	err := psql.Select("team1ID", "team2ID").
		From("games").
		Where("id = ? AND leagueID = ?", gameID, leagueID).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	if err != nil {
		return -1, -1, err
	}

	return team1, team2, nil
}

func getLosingTeamID(gameID, leagueID, winnerID int) (int, error) {
	var (
		team1 int
		team2 int
	)

	err := psql.Select("team1ID", "team2ID").
		From("games").
		Where("id = ? AND leagueID = ?", gameID, leagueID).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	if err != nil {
		return -1, err
	}

	if winnerID == team1 {
		return team2, nil
	} else if winnerID == team2 {
		return team1, nil
	} else {
		println("the winner was not either of the IDs")
		return -1, errors.New("")
	}
}

func (d *PgGamesDAO) CreateGame(leagueID, team1ID, team2ID, gameTime int) (int, error) {
	var gameID int
	err := psql.Insert("games").
		Columns("leagueID", "team1ID", "team2ID", "gametime", "complete", "winnerID", "scoreteam1", "scoreteam2").
		Values(leagueID, team1ID, team2ID, gameTime, false, -1, 0, 0).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&gameID)
	if err != nil {
		return -1, err
	}
	return gameID, nil
}

func (d *PgGamesDAO) DoesExistConflict(team1ID, team2ID, gameTime int) (bool, error) {
	//check if any game of each team is within the threshold of another scheduled game
	team1Games, err := getGamesOfTeam(team1ID)
	if err != nil {
		return false, err
	}

	for _, game := range team1Games {
		if math.Abs(float64(gameTime) - float64(game.GameTime)) < SCHEDULING_CONFLICT_THRESHOLD_SECONDS {
			return true, nil
		}
	}

	team2Games, err := getGamesOfTeam(team2ID)
	if err != nil {
		return false, err
	}

	for _, game := range team2Games {
		if math.Abs(float64(gameTime) - float64(game.GameTime)) < SCHEDULING_CONFLICT_THRESHOLD_SECONDS {
			return true, nil
		}
	}

	return false, nil
}

func (d *PgGamesDAO) GetGameInformation(gameID, leagueID int) (*GameInformation, error) {
	var gameInformation GameInformation

	err := psql.Select("*").
		From("games").
		Where("id = ? AND leagueID = ?", gameID, leagueID).
		RunWith(db).QueryRow().
		Scan(&gameInformation.Id, &gameInformation.LeagueID, &gameInformation.Team1ID, &gameInformation.Team2ID,
			&gameInformation.GameTime, &gameInformation.Complete, &gameInformation.WinnerID,
			&gameInformation.ScoreTeam1, &gameInformation.ScoreTeam2)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &gameInformation, nil
}

func (d *PgGamesDAO) HasReportResultPermissions(leagueID, gameID, userID int) (bool, error) {
	//check if user has league editResults permission
	var canReport bool
	err := psql.Select("editResults").
		From("leaguePermissions").
		Where("userID = ? AND leagueID = ?", userID, leagueID).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	//check if user has team reportResult permissions on one of the two teams
	team1ID, team2ID, err := getTeamsInGame(gameID, leagueID)

	//check for team 1
	err = psql.Select("reportResult").
		From("teamPermissions").
		Where("userID = ? AND teamID = ?", userID, team1ID).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	//check for team 2
	err = psql.Select("reportResult").
		From("teamPermissions").
		Where("userID = ? AND teamID = ?", userID, team2ID).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	return false, nil
}

func (d *PgGamesDAO) ReportGame(gameID, leagueID, winnerID, scoreTeam1, scoreTeam2 int) error {
	//update game entry
	_, err := psql.Update("games").
		Set("complete", true).
		Set("winnerID", winnerID).
		Set("scoreteam1", scoreTeam1).
		Set("scoreteam2", scoreTeam2).
		Where("id = ? AND leagueID = ?", gameID, leagueID).RunWith(db).Exec()
	if err != nil {
		return err
	}

	//update wins and losses of both teams
	_, err = db.Exec(
	`
		UPDATE teams SET wins = wins + 1
		WHERE teams.id = $1
		`, winnerID)
	if err != nil {
		return err
	}

	loserID, err := getLosingTeamID(gameID, leagueID, winnerID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
		UPDATE teams SET losses = losses + 1
		WHERE teams.id = $1
		`, loserID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}