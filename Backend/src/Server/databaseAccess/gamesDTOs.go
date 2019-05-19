package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type GameDTO struct {
	Id         int    `json:"id"`
	ExternalId string `json:"externalId"`
	LeagueId   int    `json:"leagueId"`
	Team1Id    int    `json:"team1Id"`
	Team2Id    int    `json:"team2Id"`
	GameTime   int    `json:"gameTime"`
	Complete   bool   `json:"complete"`
	WinnerId   int    `json:"winnerId"`
	LoserId    int    `json:"loserId"`
	ScoreTeam1 int    `json:"scoreTeam1"`
	ScoreTeam2 int    `json:"scoreTeam2"`
}

type GameDTOArray struct {
	rows []*GameDTO
}

func GetScannedGameDTO(rows squirrel.RowScanner) (*GameDTO, error) {
	var game GameDTO
	if err := rows.Scan(
		&game.Id,
		&game.ExternalId,
		&game.LeagueId,
		&game.Team1Id,
		&game.Team2Id,
		&game.GameTime,
		&game.Complete,
		&game.WinnerId,
		&game.LoserId,
		&game.ScoreTeam1,
		&game.ScoreTeam2,
	); err != nil {
		return nil, err
	} else {
		return &game, nil
	}
}

func (r *GameDTOArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedGameDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}
