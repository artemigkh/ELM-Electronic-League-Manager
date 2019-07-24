package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"fmt"
	"github.com/Pallinder/go-randomdata"
)

type TeamSqlDao struct{}

func tryGetUniqueIcon(leagueId int) (string, string, error) {
	// get list of icons used
	rows, err := psql.Select("icon_small").
		From("team").
		Where("league_id = ?", leagueId).
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

// Teams

//TODO: change permissions so that this can also be done during signup with giving create team perms
func (d *TeamSqlDao) CreateTeam(leagueId, userId int, teamInfo dataModel.TeamCore) (int, error) {
	iconSmall, iconLarge, err := tryGetUniqueIcon(leagueId)
	if err != nil {
		return -1, err
	}

	return d.CreateTeamWithIcon(leagueId, userId, teamInfo, iconSmall, iconLarge)
}

func (d *TeamSqlDao) CreateTeamWithIcon(leagueId, userId int, teamInfo dataModel.TeamCore,
	iconSmall, iconLarge string) (int, error) {
	var teamId = -1
	err := db.QueryRow("SELECT create_team($1,$2,$3,$4,$5,$6,$7)",
		leagueId,
		teamInfo.Name,
		teamInfo.Tag,
		teamInfo.Description,
		iconSmall,
		iconLarge,
		userId,
	).Scan(&teamId)

	return teamId, err
}

func (d *TeamSqlDao) DeleteTeam(teamId int) error {
	//remove players from team
	_, err := psql.Delete("player").
		Where("team_id = ?", teamId).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	//remove team
	_, err = psql.Delete("team").
		Where("team_id = ?", teamId).
		RunWith(db).Exec()
	return err
}

func (d *TeamSqlDao) UpdateTeam(teamId int, teamInformation dataModel.TeamCore) error {
	_, err := psql.Update("team").
		Set("name", teamInformation.Name).
		Set("tag", teamInformation.Tag).
		Set("description", teamInformation.Description).
		Where("team_id = ?", teamId).
		RunWith(db).Exec()

	return err
}

func (d *TeamSqlDao) UpdateTeamIcon(teamId int, small, large string) error {
	_, err := psql.Update("team").
		Set("icon_small", small).
		Set("icon_large", large).
		Where("team_id = ?", teamId).
		RunWith(db).Exec()

	return err
}

func (d *TeamSqlDao) GetTeamInformation(teamId int) (*dataModel.TeamWithPlayers, error) {
	rows, err := getTeamWithPlayersSelector().
		Where("team.team_id = ?", teamId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedTeamWithPlayers(rows)
}

func (d *TeamSqlDao) GetTeamWithRosters(teamId int) (*dataModel.TeamWithRosters, error) {
	rows, err := getTeamWithPlayersSelector().
		Where("team.team_id = ?", teamId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedTeamWithRosters(rows)
}

func (d *TeamSqlDao) GetAllTeamsInLeague(leagueId int) ([]*dataModel.TeamWithPlayers, error) {
	rows, err := getTeamWithPlayersSelector().
		Where("team.league_id = ?", leagueId).
		OrderBy("team.wins DESC, team.losses ASC").
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedAllTeamWithPlayers(rows)
}

func (d *TeamSqlDao) GetAllTeamsInLeagueWithRosters(leagueId int) ([]*dataModel.TeamWithRosters, error) {
	rows, err := getTeamWithPlayersSelector().
		Where("team.league_id = ?", leagueId).
		OrderBy("team.wins DESC, team.losses ASC").
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedAllTeamWithRosters(rows)
}

func (d *TeamSqlDao) GetAllTeamDisplaysInLeague(leagueId int) ([]*dataModel.TeamDisplay, error) {
	teams := TeamDisplayArray{rows: make([]*dataModel.TeamDisplay, 0)}
	if err := ScanRows(getTeamDisplaySelector().
		Where("league_id = ?", leagueId), &teams); err != nil {
		return nil, err
	}

	return teams.rows, nil
}

// Players
func (d *TeamSqlDao) CreatePlayer(leagueId, teamId int, playerInfo dataModel.PlayerCore) (int, error) {
	var playerId int
	if err := psql.Insert("player").
		Columns(
			"team_id",
			"league_id",
			"game_identifier",
			"name",
			"main_roster",
		).
		Values(
			teamId,
			leagueId,
			playerInfo.GameIdentifier,
			playerInfo.Name,
			playerInfo.MainRoster,
		).
		Suffix("RETURNING \"player_id\"").
		RunWith(db).QueryRow().Scan(&playerId); err != nil {
		return -1, err
	}

	return playerId, nil
}

func (d *TeamSqlDao) DeletePlayer(playerId int) error {
	_, err := psql.Delete("player").
		Where("player_id = ?", playerId).
		RunWith(db).Exec()
	return err
}

func (d *TeamSqlDao) UpdatePlayer(playerId int, playerInfo dataModel.PlayerCore) error {
	_, err := psql.Update("player").
		Set("game_identifier", playerInfo.GameIdentifier).
		Set("name", playerInfo.Name).
		Set("main_roster", playerInfo.MainRoster).
		Where("player_id = ?", playerId).
		RunWith(db).Exec()

	return err
}

// Get Information For Team and Player Management

func (d *TeamSqlDao) GetTeamPermissions(teamId, userId int) (*dataModel.TeamPermissionsCore, error) {
	teamPermissions, err := GetScannedTeamPermissionsCore(psql.Select(
		"administrator",
		"information",
		"games").
		From("team_permissions").
		Where("user_id = ? AND team_id = ?", userId, teamId).
		RunWith(db).QueryRow())

	if err == sql.ErrNoRows {
		return &dataModel.TeamPermissionsCore{
			Administrator: false,
			Information:   false,
			Games:         false,
		}, nil
	} else if err != nil {
		return nil, err
	}

	return teamPermissions, nil
}

func (d *TeamSqlDao) IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error) {
	var duplicateCount int

	if err := psql.Select("count(*)").
		From("team").
		Where("league_id = ? AND team_id != ? AND (name = ? OR tag = ?)", leagueId, teamId, name, tag).
		RunWith(db).QueryRow().Scan(&duplicateCount); err != nil {
		return false, "", err
	} else if duplicateCount > 0 {
		return true, "Name or tag in use", nil
	} else {
		return false, "", nil
	}
}

func (d *TeamSqlDao) DoesTeamExistInLeague(leagueId, teamId int) (bool, error) {
	var count int
	if err := psql.Select("count(*)").
		From("team").
		Where("league_id = ? AND team_id = ?", leagueId, teamId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (d *TeamSqlDao) IsTeamActive(leagueId, teamId int) (bool, error) {
	var count int
	if err := psql.Select("count(id)").
		From("game").
		Where("league_id = ? AND ( team1_id = ? OR team2_id = ?)", leagueId, teamId, teamId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (d *TeamSqlDao) DoesPlayerExist(leagueId, teamId, playerId int) (bool, error) {
	var count int
	if err := psql.Select("count(*)").
		From("player").
		Where("league_id = ? AND team_id = ? AND player_id = ?", leagueId, teamId, playerId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

// Managers

func (d *TeamSqlDao) ChangeManagerPermissions(teamId, userId int, teamPermissionInformation dataModel.TeamPermissionsCore) error {
	_, err := psql.Update("team_permissions").
		Set("administrator", teamPermissionInformation.Administrator).
		Set("information", teamPermissionInformation.Information).
		Set("games", teamPermissionInformation.Games).
		Where("team_id = ? AND user_id = ?", teamId, userId).
		RunWith(db).Exec()

	return err
}
