package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"database/sql"
)

type LeagueInformation struct {
	Id int `json:"id"`
}

type PgLeaguesDAO struct {
	psql squirrel.StatementBuilderType
}

func (d *PgLeaguesDAO) CreateLeague(userID int, name string, publicView, publicJoin bool) (int, error) {
	var leagueID int
	err := d.psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
		Values(name, publicView, publicJoin).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueID)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user ID as the league creator
	_, err = d.psql.Insert("leaguePermissions").Columns("userID", "leagueID", "editPermissions", "editTeams",
				"editUsers", "editSchedule", "editResults").Values(userID, leagueID, true, true, true, true, true).
				RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return leagueID, nil
}

func (d *PgLeaguesDAO) IsNameInUse(name string) (bool, error) {
	err := d.psql.Select("name").
		From("leagues").
		Where("name = ?", name).
		RunWith(db).QueryRow().Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
	return false, nil
}

func (d *PgLeaguesDAO) IsLeagueViewable(leagueID, userID int) (bool, error) {
	//check if publicly viewable
	var publicView bool
	err := d.psql.Select("publicview").
		From("leagues").
		Where("id = ?", leagueID).
		RunWith(db).QueryRow().Scan(&publicView)
	if err != nil {
		return false, err
	}

	if publicView {
		return true, nil
	}

	//if not publicly viewable, see if user has permission to view it. This is checked by seeing if there is a
	//leaguePermissions row with that userId and leagueId, if there is they have at least the base (viewing) privileges
	var uid int
	err = d.psql.Select("userID").
		From("leaguePermissions").
		Where("userID = ? AND leagueID = ?", userID, leagueID).
		RunWith(db).QueryRow().Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (d *PgLeaguesDAO) GetLeagueInformation(leagueID int) (*LeagueInformation, error) {
	return &LeagueInformation{Id: leagueID}, nil
}

func (d *PgLeaguesDAO) HasEditTeamsPermission(leagueID, userID int) (bool, error) {
	var canEdit bool
	err := d.psql.Select("editPermissions").
		From("leaguePermissions").
		Where("userID = ? AND leagueID = ?", userID, leagueID).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	return canEdit, nil
}