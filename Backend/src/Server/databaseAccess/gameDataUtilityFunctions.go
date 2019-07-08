package databaseAccess

import (
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
	rows []*Game
}

func GetScannedGame(rows squirrel.RowScanner) (*Game, error) {
	var (
		game  Game
		team1 TeamDisplay
		team2 TeamDisplay
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

// GameCreationInformation
func (game *GameCreationInformation) Validate(leagueId int) (bool, string, error) {
	return validate(
		game.differentTeams(),
		game.teamsExist(leagueId),
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

func (game *GameCreationInformation) teamsExist(leagueId int) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		team1Exists, team1QueryErr := teamsDAO.DoesTeamExistInLeague(leagueId, game.Team1Id)
		team2Exists, team2QueryErr := teamsDAO.DoesTeamExistInLeague(leagueId, game.Team2Id)
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

// GameTime
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

// GameResult
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
