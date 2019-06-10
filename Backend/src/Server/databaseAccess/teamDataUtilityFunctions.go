package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// TeamCore
func (team *TeamCore) validate(leagueId, teamId int) (bool, string, error) {
	return validate(
		team.name(),
		team.uniqueness(leagueId, teamId),
		team.tag())
}

func (team *TeamCore) ValidateNew(leagueId int) (bool, string, error) {
	return team.validate(leagueId, 0)
}

func (team *TeamCore) ValidateEdit(leagueId, teamId int) (bool, string, error) {
	return team.validate(leagueId, teamId)
}

func (team *TeamCore) name() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(team.Name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(team.Name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (team *TeamCore) uniqueness(leagueId, teamId int) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, problem, err := teamsDAO.IsInfoInUse(leagueId, teamId, team.Name, team.Tag)
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = problem
		} else {
			valid = true
		}
		return valid
	}
}

func (team *TeamCore) tag() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(team.Tag) > MaxTagLength {
			*problemDest = TagTooLong
		} else if len(team.Tag) < MinInformationLength {
			*problemDest = TagTooShort
		} else {
			valid = true
		}
		return valid
	}
}

// TeamPermissionsCore
type TeamPermissionsCoreArray struct {
	rows []*TeamPermissionsCore
}

func GetScannedTeamPermissionsCore(rows squirrel.RowScanner) (*TeamPermissionsCore, error) {
	var teamPermissions TeamPermissionsCore
	if err := rows.Scan(
		&teamPermissions.Administrator,
		&teamPermissions.Information,
		&teamPermissions.Players,
		&teamPermissions.ReportResults,
	); err != nil {
		return nil, err
	} else {
		return &teamPermissions, nil
	}
}

func (r *TeamPermissionsCoreArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamPermissionsCore(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}
