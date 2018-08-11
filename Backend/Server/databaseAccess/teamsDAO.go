package databaseAccess

import (
	"database/sql"
	"strings"
)

type UserInformation struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type TeamInformation struct {
	Name    string            `json:"name"`
	Tag     string            `json:"tag"`
	Wins    int               `json:"wins"`
	Losses  int               `json:"losses"`
	Members []UserInformation `json:"members"`
}

type PlayerInformation struct {
	Name string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	MainRoster bool `json:"mainRoster"`
}

type PgTeamsDAO struct{}

func (d *PgTeamsDAO) CreateTeam(leagueId, userId int, name, tag string) (int, error) {
	var teamId int
	err := psql.Insert("teams").Columns("leagueId", "name", "tag", "wins", "losses").
		Values(leagueId, name, strings.ToUpper(tag), 0, 0).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&teamId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("teamPermissions").
		Columns("userId", "teamId", "editPermissions", "editTeamInfo", "editUsers", "reportResult").
		Values(userId, teamId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return teamId, nil
}

func (d *PgTeamsDAO) IsInfoInUse(name, tag string, leagueId int) (bool, string, error) {
	//check if name in use
	err := psql.Select("name").
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
	err = psql.Select("tag").
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

func (d *PgTeamsDAO) GetTeamInformation(teamId, leagueId int) (*TeamInformation, error) {
	var teamInformation TeamInformation
	//get team information
	err := psql.Select("name", "tag", "wins", "losses").
		From("teams").
		Where("id = ? AND leagueId = ?", teamId, leagueId).
		RunWith(db).QueryRow().Scan(&teamInformation.Name, &teamInformation.Tag, &teamInformation.Wins, &teamInformation.Losses)
	if err != nil {
		return nil, err
	}

	//get users of team
	var members []UserInformation

	rows, err := db.Query(`
		SELECT id, email FROM users
			WHERE id IN
				(
					SELECT userId FROM teamPermissions
					WHERE teamId = $1
				)
	`, teamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var member UserInformation
	for rows.Next() {
		err := rows.Scan(&member.Id, &member.Email)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	teamInformation.Members = members
	return &teamInformation, nil
}

func (d *PgTeamsDAO) DoesTeamExist(teamId, leagueId int) (bool, error) {
	var name string
	err := psql.Select("name").
		From("teams").
		Where("id = ? AND leagueId = ?", teamId, leagueId).
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
