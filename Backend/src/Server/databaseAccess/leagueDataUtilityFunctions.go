package databaseAccess

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

// League
type LeagueArray struct {
	rows []*League
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

func GetScannedLeague(rows squirrel.RowScanner) (*League, error) {
	var league League
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

// LeagueCore
func (league *LeagueCore) validate(leagueId int) (bool, string, error) {
	fmt.Printf("league is %+v\n", league)
	return validate(
		league.name(),
		league.uniqueness(leagueId),
		league.description(),
		league.game(),
		league.permissions(),
		league.timestamps())
}

func (league *LeagueCore) ValidateNew() (bool, string, error) {
	return league.validate(0)
}

func (league *LeagueCore) ValidateEdit(leagueId int) (bool, string, error) {
	return league.validate(leagueId)
}

func (league *LeagueCore) name() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(league.Name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(league.Name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (league *LeagueCore) uniqueness(leagueId int) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, err := Leagues.IsNameInUse(leagueId, league.Name)
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = LeagueNameInUse
		} else {
			valid = true
		}
		return valid
	}
}

func (league *LeagueCore) description() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if len(league.Description) > MaxDescriptionLength {
			*problemDest = DescriptionTooLong
			return false
		} else {
			return true
		}
	}
}

func (league *LeagueCore) game() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		for _, g := range ValidGameStrings {
			if g == league.Game {
				return true
			}
		}
		*problemDest = InvalidGame
		return false
	}
}

func (league *LeagueCore) permissions() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if league.PublicJoin && !league.PublicView {
			*problemDest = LeaguePermissionsWrong
			return false
		} else {
			return true
		}
	}
}

func (league *LeagueCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if league.SignupStart > league.SignupEnd ||
			league.LeagueStart > league.LeagueEnd {
			*problemDest = TimeOutOfOrder
		} else if league.LeagueStart < league.SignupEnd {
			*problemDest = PeriodOutOfOrder
		} else {
			valid = true
		}
		return valid
	}
}

// Markdown
func (md *Markdown) Validate() (bool, string, error) {
	return validate(md.markdown())
}

func (md *Markdown) markdown() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if len(md.Markdown) > MaxMdLength {
			*problemDest = MarkdownTooLong
			return false
		} else {
			return true
		}
	}
}

// LeaguePermissionsCore
func GetScannedLeaguePermissionsCore(rows squirrel.RowScanner) (*LeaguePermissionsCore, error) {
	var leaguePermissions LeaguePermissionsCore
	if err := rows.Scan(
		&leaguePermissions.Administrator,
		&leaguePermissions.CreateTeams,
		&leaguePermissions.EditTeams,
		&leaguePermissions.EditGames,
	); err != nil {
		return nil, err
	} else {
		return &leaguePermissions, nil
	}
}

func (p *LeaguePermissionsCore) Validate() (bool, string, error) {
	return validate(p.consistent())
}

func (p *LeaguePermissionsCore) consistent() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if (p.CreateTeams || p.EditGames || p.EditTeams) && p.Administrator {
			*problemDest = AdminLackingPermissions
			return false
		} else {
			return true
		}
	}
}
