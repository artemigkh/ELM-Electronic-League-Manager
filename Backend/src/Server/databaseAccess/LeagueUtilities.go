package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// League
type LeagueArray struct {
	rows []*dataModel.League
}

func getLeagueSelector() squirrel.SelectBuilder {
	return psql.Select(
		"league_id",
		"name",
		"description",
		"game",
		"public_view",
		"public_join",
		"signup_start",
		"signup_end",
		"league_start",
		"league_end",
	).From("league")
}

func GetScannedLeague(rows squirrel.RowScanner) (*dataModel.League, error) {
	var league dataModel.League
	if err := rows.Scan(
		&league.LeagueId,
		&league.Name,
		&league.Description,
		&league.Game,
		&league.PublicView,
		&league.PublicJoin,
		&league.SignupStart,
		&league.SignupEnd,
		&league.LeagueStart,
		&league.LeagueEnd,
	); err != nil {
		return nil, err
	} else {
		return &league, nil
	}
}

func (r *LeagueArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedLeague(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

// LeaguePermissionsCore
func GetScannedLeaguePermissionsCore(rows squirrel.RowScanner) (*dataModel.LeaguePermissionsCore, error) {
	var leaguePermissions dataModel.LeaguePermissionsCore
	if err := rows.Scan(
		&leaguePermissions.Administrator,
		&leaguePermissions.CreateTeams,
		&leaguePermissions.EditTeams,
		&leaguePermissions.EditGames,
	); err != nil {
		if err == sql.ErrNoRows {
			return &dataModel.LeaguePermissionsCore{
				Administrator: false,
				CreateTeams:   false,
				EditTeams:     false,
				EditGames:     false,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return &leaguePermissions, nil
	}
}
