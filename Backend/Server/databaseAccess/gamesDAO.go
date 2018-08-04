package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"math"
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

type PgGamesDAO struct {
	psql squirrel.StatementBuilderType
}

func getGamesOfTeam(psql squirrel.StatementBuilderType, teamId int) ([]GameInformation, error) {
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

func (d *PgGamesDAO) CreateGame(leagueID, team1ID, team2ID, gameTime int) (int, error) {
	var gameID int
	err := d.psql.Insert("games").
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
	team1Games, err := getGamesOfTeam(d.psql, team1ID)
	if err != nil {
		return false, err
	}

	for _, game := range team1Games {
		if math.Abs(float64(gameTime) - float64(game.GameTime)) < SCHEDULING_CONFLICT_THRESHOLD_SECONDS {
			return true, nil
		}
	}

	team2Games, err := getGamesOfTeam(d.psql, team2ID)
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

func (d *PgGamesDAO) GetGameInformation(gameID int) (*GameInformation, error) {
	var gameInformation GameInformation

	err := d.psql.Select("*").
		From("games").
		Where("id = ?", gameID).
		RunWith(db).QueryRow().
		Scan(&gameInformation.Id, &gameInformation.LeagueID, &gameInformation.Team1ID, &gameInformation.Team2ID,
			&gameInformation.GameTime, &gameInformation.Complete, &gameInformation.WinnerID,
			&gameInformation.ScoreTeam1, &gameInformation.ScoreTeam2)
	if err != nil {
		return nil, err
	}

	return &gameInformation, nil
}
