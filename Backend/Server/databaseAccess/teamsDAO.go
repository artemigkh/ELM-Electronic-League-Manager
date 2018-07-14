package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"database/sql"
	"strings"
)

type PgTeamsDAO struct {
	psql squirrel.StatementBuilderType
}

func (d* PgTeamsDAO) CreateTeam(leagueID, userID int, name, tag string) (int, error) {
	var teamID int
	err := d.psql.Insert("teams").Columns("leagueID", "name", "tag", "wins", "losses").
		Values(leagueID, name, strings.ToUpper(tag), 0, 0).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&teamID)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user ID as the league creator
	_, err = d.psql.Insert("teamPermissions").
		Columns("userID", "teamID", "editPermissions", "editTeamInfo", "editUsers", "reportResult").
		Values(userID, teamID, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return teamID, nil
}

func (d* PgTeamsDAO)IsInfoInUse(name, tag string, leagueID int) (bool, string, error) {
	//check if name in use
	err := d.psql.Select("name").
		From("teams").
		Where("name = ?", name).
		RunWith(db).QueryRow().Scan(&name)
	if err == sql.ErrNoRows {
		//check for tag
	} else if err != nil {
		return false, "", err
	} else {
		return true, "nameInUse", nil
	}

	//check if name in use
	err = d.psql.Select("tag").
		From("teams").
		Where("tag = ?", strings.ToUpper(tag)).
		RunWith(db).QueryRow().Scan(&tag)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	} else {
		return true, "tagInUse", nil
	}
}