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

func (d *PgLeaguesDAO) CreateLeague(userId int, name string, publicView, publicJoin bool) (int, error) {
	var leagueId int
	err := psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
		Values(name, publicView, publicJoin).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("leaguePermissions").
		Columns("userId", "leagueId", "editPermissions", "createTeams",
			"editTeams", "editUsers", "editSchedule", "editResults").
		Values(userId, leagueId, true, true, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return leagueId, nil
}

func (d *PgLeaguesDAO) JoinLeague(userId, leagueId int) error {
	_, err := psql.Insert("leaguePermissions").
		Columns("userId", "leagueId", "editPermissions", "createTeams",
			"editTeams", "editUsers", "editSchedule", "editResults").
		Values(userId, leagueId, false, true, false, false, false, false).
		RunWith(db).Exec()
	if err != nil {
		return err
	} else {
		return nil
	}
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

func (d *PgLeaguesDAO) IsLeagueViewable(leagueId, userId int) (bool, error) {
	//check if publicly viewable
	var publicView bool
	err := psql.Select("publicview").
		From("leagues").
		Where("id = ?", leagueId).
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
	err = psql.Select("userId").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (d *PgLeaguesDAO) GetLeagueInformation(leagueId int) (*LeagueInformation, error) {
	return &LeagueInformation{Id: leagueId}, nil
}

func (d *PgLeaguesDAO) HasEditTeamsPermission(leagueId, userId int) (bool, error) {
	var canEdit bool
	err := psql.Select("editPermissions").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	return canEdit, nil
}

func (d *PgLeaguesDAO) HasCreateTeamsPermission(leagueId, userId int) (bool, error) {
	var canEdit bool
	err := psql.Select("createTeams").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	return canEdit, nil
}

func (d *PgLeaguesDAO) GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error) {
	rows, err := psql.Select("id", "name", "tag", "wins", "losses").From("teams").
		Where("leagueId = ?", leagueId).
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

//TODO: make invite system for private leagues, check if user invited in this function
func (d *PgLeaguesDAO) CanJoinLeague(userId, leagueId int) (bool, error) {
	var canJoin bool
	err := psql.Select("publicJoin").
		From("leagues").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&canJoin)
	if err != nil {
		return false, err
	}

	return canJoin, nil
}
