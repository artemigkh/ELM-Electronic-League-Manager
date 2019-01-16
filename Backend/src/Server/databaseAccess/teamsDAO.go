package databaseAccess

import (
	"database/sql"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"strings"
)

type TeamInformation struct {
	Name        string              `json:"name"`
	Tag         string              `json:"tag"`
	Description string              `json:"description"`
	Wins        int                 `json:"wins"`
	Losses      int                 `json:"losses"`
	IconSmall   string              `json:"iconSmall"`
	IconLarge   string              `json:"iconLarge"`
	Players     []PlayerInformation `json:"players"`
}

type TeamPermissions struct {
	Administrator bool
	Information   bool
	Players       bool
	ReportResults bool
}

type PlayerInformation struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	MainRoster     bool   `json:"mainRoster"`
}

type PgTeamsDAO struct{}

func tryGetUniqueIcon(leagueId int) (string, string, error) {
	// get list of icons used
	rows, err := psql.Select("iconSmall").
		From("teams").
		Where("leagueId = ?", leagueId).
		RunWith(db).Query()

	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	// generate bool who's indices indicate if that number is available
	var availableIcons []bool
	for i := 0; i < 9; i++ {
		availableIcons = append(availableIcons, true)
	}

	// mark numbers as taken if the filename associated with it is present
	var icon string
	for rows.Next() {
		err := rows.Scan(&icon)
		if err != nil {
			return "", "", err
		}
		for i := 0; i < 9; i++ {
			if icon == fmt.Sprintf("generic-%v-small.png", i+1) {
				availableIcons[i] = false
			}
		}
	}
	if rows.Err() != nil {
		return "", "", err
	}

	// create list of available generic icons
	var availableNumbers []int
	for i := 0; i < 9; i++ {
		if availableIcons[i] {
			availableNumbers = append(availableNumbers, i+1)
		}
	}

	// select one either from available or if all taken a random one
	var newIconNumber int
	println(fmt.Sprintf("Available numbers: %v", len(availableNumbers)))
	if len(availableNumbers) == 0 {
		newIconNumber = randomdata.Number(1, 9)
	} else if len(availableNumbers) == 1 {
		newIconNumber = availableNumbers[0]
	} else {
		newIconNumber = availableNumbers[randomdata.Number(0, len(availableNumbers)-1)]
	}

	return fmt.Sprintf("generic-%v-small.png", newIconNumber),
		fmt.Sprintf("generic-%v-large.png", newIconNumber), nil
}

func (d *PgTeamsDAO) CreateTeam(leagueId, userId int, name, tag, description string) (int, error) {
	smallIcon, largeIcon, err := tryGetUniqueIcon(leagueId)
	if err != nil {
		return -1, err
	}

	var teamId int
	err = psql.Insert("teams").
		Columns("leagueId", "name", "tag", "description", "wins", "losses", "iconSmall", "iconLarge").
		Values(leagueId, name, strings.ToUpper(tag), description, 0, 0, smallIcon, largeIcon).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&teamId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("teamPermissions").
		Columns("userId", "teamId", "administrator", "information", "players", "reportResults").
		Values(userId, teamId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return teamId, nil
}

func (d *PgTeamsDAO) CreateTeamWithIcon(leagueId, userId int, name, tag, description, small, large string) (int, error) {
	var teamId int
	err := psql.Insert("teams").
		Columns("leagueId", "name", "tag", "description", "wins", "losses", "iconSmall", "iconLarge").
		Values(leagueId, name, strings.ToUpper(tag), description, 0, 0, small, large).Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&teamId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("teamPermissions").
		Columns("userId", "teamId", "administrator", "information", "players", "reportResults").
		Values(userId, teamId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return teamId, nil
}

func (d *PgTeamsDAO) IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error) {
	//check if name in use
	var teamIdOfMatch int
	err := psql.Select("name, id").
		From("teams").
		Where("name = ? AND leagueId = ?", name, leagueId).
		RunWith(db).QueryRow().Scan(&name, &teamIdOfMatch)
	if err == sql.ErrNoRows {
		//check for tag
	} else if err != nil {
		return false, "", err
	} else if teamId != teamIdOfMatch {
		return true, "nameInUse", nil
	}

	//check if tag in use
	err = psql.Select("tag, id").
		From("teams").
		Where("tag = ? AND leagueId = ?", strings.ToUpper(tag), leagueId).
		RunWith(db).QueryRow().Scan(&tag, &teamIdOfMatch)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	} else if teamId != teamIdOfMatch {
		return true, "tagInUse", nil
	} else {
		return false, "", nil
	}
}

func (d *PgTeamsDAO) GetTeamInformation(leagueId, teamId int) (*TeamInformation, error) {
	var teamInformation TeamInformation
	//get team information
	err := psql.Select("name", "tag", "description", "wins", "losses", "iconSmall", "iconLarge").
		From("teams").
		Where("id = ? AND leagueId = ?", teamId, leagueId).
		RunWith(db).QueryRow().Scan(&teamInformation.Name, &teamInformation.Tag, &teamInformation.Description,
		&teamInformation.Wins, &teamInformation.Losses, &teamInformation.IconSmall, &teamInformation.IconLarge)
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

func doesTeamExist(leagueId, teamId int) (bool, error) {
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
}

func (d *PgTeamsDAO) DoesTeamExist(leagueId, teamId int) (bool, error) {
	return doesTeamExist(leagueId, teamId)
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

func (d *PgTeamsDAO) UpdateTeam(leagueId, teamId int, name, tag, description string) error {
	_, err := db.Exec(
		`
		UPDATE teams SET name = $1, tag = $2, description = $3
		WHERE id = $4 AND leagueId = $5
		`, name, tag, description, teamId, leagueId)
	return err
}

func (d *PgTeamsDAO) RemovePlayer(teamId, playerId int) error {
	_, err := psql.Delete("players").
		Where("id = ? AND teamId = ?", playerId, teamId).
		RunWith(db).Exec()
	return err
}

func (d *PgTeamsDAO) UpdatePlayer(teamId, playerId int, gameIdentifier, name string, mainRoster bool) error {
	_, err := db.Exec(
		`
		UPDATE players SET gameIdentifier = $1, name = $2, mainRoster = $3
		WHERE id = $4 AND teamId = $5
		`, gameIdentifier, name, mainRoster, playerId, teamId)
	return err
}

func (d *PgTeamsDAO) DoesPlayerExist(teamId, playerId int) (bool, error) {
	var name string
	err := psql.Select("name").
		From("players").
		Where("id = ? AND teamId = ?", playerId, teamId).
		RunWith(db).QueryRow().Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *PgTeamsDAO) IsTeamActive(leagueId, teamId int) (bool, error) {
	var gameId int
	err := psql.Select("id").
		From("games").
		Where("leagueId = ? AND ( team1Id = ? OR team2Id = ?)", leagueId, teamId, teamId).
		RunWith(db).QueryRow().Scan(&gameId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (d *PgTeamsDAO) DeleteTeam(leagueId, teamId int) error {
	//remove players from team
	_, err := psql.Delete("players").
		Where("teamId = ?", teamId).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	//remove team
	_, err = psql.Delete("teams").
		Where("id = ? AND leagueId = ?", teamId, leagueId).
		RunWith(db).Exec()
	return err
}

func (d *PgTeamsDAO) GetTeamPermissions(teamId, userId int) (*TeamPermissions, error) {
	var tp TeamPermissions
	err := psql.Select("administrator", "information", "players", "reportResults").
		From("teamPermissions").
		Where("userId = ? AND teamId = ?", userId, teamId).
		RunWith(db).QueryRow().Scan(&tp.Administrator, &tp.Information, &tp.Players, &tp.ReportResults)
	if err == sql.ErrNoRows {
		return &TeamPermissions{
			Administrator: false,
			Information:   false,
			Players:       false,
			ReportResults: false,
		}, nil
	} else if err != nil {
		return nil, err
	}

	return &tp, nil
}

func (d *PgTeamsDAO) ChangeManagerPermissions(teamId, userId int, administrator, information, players, reportResults bool) error {
	_, err := db.Exec(
		`
		UPDATE teamPermissions SET administrator = $1, information = $2, players = $3, reportResults = $4
		WHERE teamId = $5 AND userId = $6
		`, administrator, information, players, reportResults, teamId, userId)
	return err
}
