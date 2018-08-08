package databaseAccess

import (
	"database/sql"
	"strings"
)

type UserInformation struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

type TeamInformation struct {
	Name string `json:"name"`
	Tag string `json:"tag"`
	Wins int `json:"wins"`
	Losses int `json:"losses"`
	Members []UserInformation `json:"members"`
}

type PgTeamsDAO struct {}

func (d* PgTeamsDAO) CreateTeam(leagueID, userID int, name, tag string) (int, error) {
	var teamID int
	err := psql.Insert("teams").Columns("leagueID", "name", "tag", "wins", "losses").
		Values(leagueID, name, strings.ToUpper(tag), 0, 0).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&teamID)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user ID as the league creator
	_, err = psql.Insert("teamPermissions").
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

func (d *PgTeamsDAO) GetTeamInformation(teamID, leagueID int) (*TeamInformation, error) {
	var teamInformation TeamInformation
	//get team information
	err := psql.Select("name", "tag", "wins", "losses").
		From("teams").
		Where("id = ? AND leagueID = ?", teamID, leagueID).
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
					SELECT userID FROM teamPermissions
					WHERE teamID = $1
				)
	`, teamID)
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

func (d *PgTeamsDAO) DoesTeamExist(teamID, leagueID int) (bool, error) {
	var name string
	err := psql.Select("name").
		From("teams").
		Where("id = ? AND leagueID = ?", teamID, leagueID).
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