package databaseAccess

import (
	"database/sql"
)

type LeagueInformation struct {
	Id int `json:"id"`
}

type TeamSummaryInformation struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}

type PgLeaguesDAO struct{}

func (d *PgLeaguesDAO) CreateLeague(userID int, name string, publicView, publicJoin bool) (int, error) {
	var leagueID int
	err := psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
		Values(name, publicView, publicJoin).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueID)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user ID as the league creator
	_, err = psql.Insert("leaguePermissions").Columns("userID", "leagueID", "editPermissions", "editTeams",
		"editUsers", "editSchedule", "editResults").Values(userID, leagueID, true, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return leagueID, nil
}

func (d *PgLeaguesDAO) IsNameInUse(name string) (bool, error) {
	err := psql.Select("name").
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
	err := psql.Select("publicview").
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
	err = psql.Select("userID").
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
	err := psql.Select("editPermissions").
		From("leaguePermissions").
		Where("userID = ? AND leagueID = ?", userID, leagueID).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	return canEdit, nil
}

func (d *PgLeaguesDAO) GetTeamSummary(leagueID int) ([]TeamSummaryInformation, error) {
	rows, err := psql.Select("id", "name", "tag", "wins", "losses").From("teams").
		Where("leagueID = ?", leagueID).
		OrderBy("wins DESC").
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []TeamSummaryInformation
	var team TeamSummaryInformation

	for rows.Next() {
		err := rows.Scan(&team.Id, &team.Name, &team.Tag, &team.Wins, &team.Losses)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return teams, nil
}
