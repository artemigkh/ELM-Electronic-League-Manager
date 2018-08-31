package databaseAccess

import (
	"database/sql"
	"strings"
)

type TeamInformation struct {
	Name    string              `json:"name"`
	Tag     string              `json:"tag"`
	Wins    int                 `json:"wins"`
	Losses  int                 `json:"losses"`
	Players []PlayerInformation `json:"players"`
}

type PlayerInformation struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	MainRoster     bool   `json:"mainRoster"`
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
		Columns("userId", "teamId", "editPermissions", "editTeamInfo", "editPlayers", "reportResult").
		Values(userId, teamId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return teamId, nil
}

func (d *PgTeamsDAO) IsInfoInUse(leagueId int, name, tag string) (bool, string, error) {
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

func (d *PgTeamsDAO) GetTeamInformation(leagueId, teamId int) (*TeamInformation, error) {
	var teamInformation TeamInformation
	//get team information
	err := psql.Select("name", "tag", "wins", "losses").
		From("teams").
		Where("id = ? AND leagueId = ?", teamId, leagueId).
		RunWith(db).QueryRow().Scan(&teamInformation.Name, &teamInformation.Tag, &teamInformation.Wins, &teamInformation.Losses)
	if err != nil {
		return nil, err
	}

	//get players of team
	rows, err := psql.Select("id", "gameIdentifier", "name", "mainRoster").
		From("players").
		Where("teamId = ?", teamId).RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []PlayerInformation
	var player PlayerInformation

	for rows.Next() {
		err := rows.Scan(&player.Id, &player.GameIdentifier, &player.Name, &player.MainRoster)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	teamInformation.Players = players
	return &teamInformation, nil
}

func (d *PgTeamsDAO) DoesTeamExist(leagueId, teamId int) (bool, error) {
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

func (d *PgTeamsDAO) HasPlayerEditPermissions(leagueId, teamId, userId int) (bool, error) {
	//check if league admin
	var canEdit bool
	err := psql.Select("editTeams").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	if canEdit {
		return true, nil
	}

	//check if team admin
	err = psql.Select("editPlayers").
		From("teamPermissions").
		Where("userId = ? AND teamId = ?", userId, teamId).
		RunWith(db).QueryRow().Scan(&canEdit)
	if err != nil {
		return false, err
	}

	if canEdit {
		return true, nil
	}

	return false, nil
}

func (d *PgTeamsDAO) AddNewPlayer(teamId int, gameIdentifier, name string, mainRoster bool) (int, error) {
	var playerId int
	err := psql.Insert("players").Columns("teamId", "gameIdentifier", "name", "mainRoster").
		Values(teamId, gameIdentifier, name, mainRoster).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&playerId)
	if err != nil {
		return -1, err
	}

	return playerId, nil
}
