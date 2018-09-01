package databaseAccess

import (
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"math"
)

const (
	SCHEDULING_CONFLICT_THRESHOLD_SECONDS = 120
)

type GameInformation struct {
	Id         int  `json:"id"`
	LeagueId   int  `json:"leagueId"`
	Team1Id    int  `json:"team1Id"`
	Team2Id    int  `json:"team2Id"`
	GameTime   int  `json:"gameTime"`
	Complete   bool `json:"complete"`
	WinnerId   int  `json:"winnerId"`
	ScoreTeam1 int  `json:"scoreTeam1"`
	ScoreTeam2 int  `json:"scoreTeam2"`
}

type PgGamesDAO struct{}

func getGamesOfTeam(teamId int) ([]GameInformation, error) {
	rows, err := psql.Select("*").From("games").
		Where(squirrel.Or{squirrel.Eq{"team1Id": teamId}, squirrel.Eq{"team2Id": teamId}}).
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameInformation
	var game GameInformation

	for rows.Next() {
		err := rows.Scan(&game.Id, &game.LeagueId, &game.Team1Id, &game.Team2Id, &game.GameTime,
			&game.Complete, &game.WinnerId, &game.ScoreTeam1, &game.ScoreTeam2)
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

func getTeamsInGame(gameId, leagueId int) (int, int, error) {
	var (
		team1 int
		team2 int
	)

	err := psql.Select("team1Id", "team2Id").
		From("games").
		Where("id = ? AND leagueId = ?", gameId, leagueId).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	if err != nil {
		return -1, -1, err
	}

	return team1, team2, nil
}

func getLosingTeamId(gameId, leagueId, winnerId int) (int, error) {
	var (
		team1 int
		team2 int
	)

	err := psql.Select("team1Id", "team2Id").
		From("games").
		Where("id = ? AND leagueId = ?", gameId, leagueId).
		RunWith(db).QueryRow().Scan(&team1, &team2)
	if err != nil {
		return -1, err
	}

	if winnerId == team1 {
		return team2, nil
	} else if winnerId == team2 {
		return team1, nil
	} else {
		println("the winner was not either of the Ids")
		return -1, errors.New("")
	}
}

func (d *PgGamesDAO) CreateGame(leagueId, team1Id, team2Id, gameTime int) (int, error) {
	var gameId int
	err := psql.Insert("games").
		Columns("leagueId", "team1Id", "team2Id", "gametime", "complete", "winnerId", "scoreteam1", "scoreteam2").
		Values(leagueId, team1Id, team2Id, gameTime, false, -1, 0, 0).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&gameId)
	if err != nil {
		return -1, err
	}
	return gameId, nil
}

func (d *PgGamesDAO) DoesExistConflict(team1Id, team2Id, gameTime int) (bool, error) {
	//check if any game of each team is within the threshold of another scheduled game
	team1Games, err := getGamesOfTeam(team1Id)
	if err != nil {
		return false, err
	}

	for _, game := range team1Games {
		if math.Abs(float64(gameTime)-float64(game.GameTime)) < SCHEDULING_CONFLICT_THRESHOLD_SECONDS {
			return true, nil
		}
	}

	team2Games, err := getGamesOfTeam(team2Id)
	if err != nil {
		return false, err
	}

	for _, game := range team2Games {
		if math.Abs(float64(gameTime)-float64(game.GameTime)) < SCHEDULING_CONFLICT_THRESHOLD_SECONDS {
			return true, nil
		}
	}

	return false, nil
}

func (d *PgGamesDAO) GetGameInformation(leagueId, gameId int) (*GameInformation, error) {
	var gameInformation GameInformation

	err := psql.Select("*").
		From("games").
		Where("id = ? AND leagueId = ?", gameId, leagueId).
		RunWith(db).QueryRow().
		Scan(&gameInformation.Id, &gameInformation.LeagueId, &gameInformation.Team1Id, &gameInformation.Team2Id,
			&gameInformation.GameTime, &gameInformation.Complete, &gameInformation.WinnerId,
			&gameInformation.ScoreTeam1, &gameInformation.ScoreTeam2)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &gameInformation, nil
}

func (d *PgGamesDAO) HasReportResultPermissions(leagueId, gameId, userId int) (bool, error) {
	//check if user has league editResults permission
	var canReport bool
	err := psql.Select("editResults").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	//check if user has team reportResult permissions on one of the two teams
	team1Id, team2Id, err := getTeamsInGame(gameId, leagueId)

	//check for team 1
	err = psql.Select("reportResult").
		From("teamPermissions").
		Where("userId = ? AND teamId = ?", userId, team1Id).
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
		Where("userId = ? AND teamId = ?", userId, team2Id).
		RunWith(db).QueryRow().Scan(&canReport)
	if err != nil {
		return false, err
	}

	if canReport {
		return true, nil
	}

	return false, nil
}

func (d *PgGamesDAO) ReportGame(leagueId, gameId, winnerId, scoreTeam1, scoreTeam2 int) error {
	//update game entry
	_, err := psql.Update("games").
		Set("complete", true).
		Set("winnerId", winnerId).
		Set("scoreteam1", scoreTeam1).
		Set("scoreteam2", scoreTeam2).
		Where("id = ? AND leagueId = ?", gameId, leagueId).RunWith(db).Exec()
	if err != nil {
		return err
	}

	//update wins and losses of both teams
	_, err = db.Exec(
		`
		UPDATE teams SET wins = wins + 1
		WHERE teams.id = $1
		`, winnerId)
	if err != nil {
		return err
	}

	loserId, err := getLosingTeamId(gameId, leagueId, winnerId)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
		UPDATE teams SET losses = losses + 1
		WHERE teams.id = $1
		`, loserId)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *PgGamesDAO) DeleteGame(leagueId, gameId int) error {
	_, err := psql.Delete("games").
		Where("id = ? AND leagueId = ?", gameId, leagueId).
		RunWith(db).Exec()
	return err
}
