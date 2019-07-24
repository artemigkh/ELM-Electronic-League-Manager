package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// Game
func getGameSelector() squirrel.SelectBuilder {
	return psql.Select(
		"game_id",
		"complete",
		"game_time",
		"winner_id",
		"loser_id",
		"score_team1",
		"score_team2",
		"team1.team_id",
		"team1.name",
		"team1.tag",
		"team1.icon_small",
		"team1.wins",
		"team1.losses",
		"team2.team_id",
		"team2.name",
		"team2.tag",
		"team2.icon_small",
		"team2.wins",
		"team2.losses").
		From("game").
		Join("team AS team1 ON game.team1_id = team1.team_id").
		Join("team AS team2 ON game.team2_id = team2.team_id")
}

type GameArray struct {
	rows []*dataModel.Game
}

func GetScannedGame(rows squirrel.RowScanner) (*dataModel.Game, error) {
	var (
		game  dataModel.Game
		team1 dataModel.TeamDisplay
		team2 dataModel.TeamDisplay
	)
	if err := rows.Scan(
		&game.GameId,
		&game.Complete,
		&game.GameTime,
		&game.WinnerId,
		&game.LoserId,
		&game.ScoreTeam1,
		&game.ScoreTeam2,
		&team1.TeamId,
		&team1.Name,
		&team1.Tag,
		&team1.IconSmall,
		&team1.Wins,
		&team1.Losses,
		&team2.TeamId,
		&team2.Name,
		&team2.Tag,
		&team2.IconSmall,
		&team2.Wins,
		&team2.Losses,
	); err != nil {
		return nil, err
	} else {
		game.Team1 = team1
		game.Team2 = team2
		return &game, nil
	}
}

func (r *GameArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedGame(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}
