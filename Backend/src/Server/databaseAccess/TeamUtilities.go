package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// TeamPermissionsCore
type TeamPermissionsCoreArray struct {
	rows []*dataModel.TeamPermissionsCore
}

func GetScannedTeamPermissionsCore(rows squirrel.RowScanner) (*dataModel.TeamPermissionsCore, error) {
	var teamPermissions dataModel.TeamPermissionsCore
	if err := rows.Scan(
		&teamPermissions.Administrator,
		&teamPermissions.Information,
		&teamPermissions.Games,
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

// TeamWithPlayers
type TeamWithPlayersArray struct {
	rows []*dataModel.TeamWithPlayers
}

func getTeamWithPlayersSelector() squirrel.SelectBuilder {
	return psql.Select(
		"team.team_id",
		"team.name",
		"team.description",
		"team.tag",
		"team.icon_small",
		"team.icon_large",
		"team.wins",
		"team.losses",
		"player.player_id",
		"player.name",
		"player.game_identifier",
		"player.main_roster",
	).
		From("team").
		LeftJoin("player ON team.team_id = player.team_id")
}

func GetScannedAllTeamWithRosters(rows *sql.Rows) ([]*dataModel.TeamWithRosters, error) {
	teams := make([]*dataModel.TeamWithRosters, 0)
	getUniqueTeam := func(newTeam *dataModel.TeamWithRosters) *dataModel.TeamWithRosters {
		for _, team := range teams {
			if newTeam.TeamId == team.TeamId {
				return team
			}
		}
		teams = append(teams, newTeam)
		return newTeam
	}

	defer rows.Close()
	for rows.Next() {
		var team dataModel.TeamWithRosters
		team.SubstituteRoster = make([]*dataModel.Player, 0)
		team.MainRoster = make([]*dataModel.Player, 0)
		var (
			playerId             sql.NullInt64
			playerName           sql.NullString
			playerGameIdentifier sql.NullString
			playerMainRoster     sql.NullBool
		)
		if err := rows.Scan(
			&team.TeamId,
			&team.Name,
			&team.Description,
			&team.Tag,
			&team.IconSmall,
			&team.IconLarge,
			&team.Wins,
			&team.Losses,
			&playerId,
			&playerName,
			&playerGameIdentifier,
			&playerMainRoster,
		); err != nil {
			return nil, err
		}

		uniqueTeam := getUniqueTeam(&team)
		if playerId.Valid {
			if playerMainRoster.Bool {
				uniqueTeam.MainRoster = append(uniqueTeam.MainRoster, &dataModel.Player{
					PlayerId:       int(playerId.Int64),
					Name:           playerName.String,
					GameIdentifier: playerGameIdentifier.String,
					MainRoster:     playerMainRoster.Bool,
				})
			} else {
				uniqueTeam.SubstituteRoster = append(uniqueTeam.SubstituteRoster, &dataModel.Player{
					PlayerId:       int(playerId.Int64),
					Name:           playerName.String,
					GameIdentifier: playerGameIdentifier.String,
					MainRoster:     playerMainRoster.Bool,
				})
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func GetScannedTeamWithRosters(rows *sql.Rows) (*dataModel.TeamWithRosters, error) {
	teams, err := GetScannedAllTeamWithRosters(rows)
	if err != nil {
		return nil, err
	} else {
		return teams[0], nil
	}
}

func GetScannedAllTeamWithPlayers(rows *sql.Rows) ([]*dataModel.TeamWithPlayers, error) {
	teams := make([]*dataModel.TeamWithPlayers, 0)
	getUniqueTeam := func(newTeam *dataModel.TeamWithPlayers) *dataModel.TeamWithPlayers {
		for _, team := range teams {
			if newTeam.TeamId == team.TeamId {
				return team
			}
		}
		teams = append(teams, newTeam)
		return newTeam
	}

	defer rows.Close()
	for rows.Next() {
		var team dataModel.TeamWithPlayers
		team.Players = make([]*dataModel.Player, 0)
		var (
			playerId             sql.NullInt64
			playerName           sql.NullString
			playerGameIdentifier sql.NullString
			playerMainRoster     sql.NullBool
		)
		if err := rows.Scan(
			&team.TeamId,
			&team.Name,
			&team.Description,
			&team.Tag,
			&team.IconSmall,
			&team.IconLarge,
			&team.Wins,
			&team.Losses,
			&playerId,
			&playerName,
			&playerGameIdentifier,
			&playerMainRoster,
		); err != nil {
			return nil, err
		}

		uniqueTeam := getUniqueTeam(&team)
		if playerId.Valid {
			uniqueTeam.Players = append(uniqueTeam.Players, &dataModel.Player{
				PlayerId:       int(playerId.Int64),
				Name:           playerName.String,
				GameIdentifier: playerGameIdentifier.String,
				MainRoster:     playerMainRoster.Bool,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func GetScannedTeamWithPlayers(rows *sql.Rows) (*dataModel.TeamWithPlayers, error) {
	teams, err := GetScannedAllTeamWithPlayers(rows)
	if err != nil {
		return nil, err
	} else {
		return teams[0], nil
	}
}

func getLoLTeamStubSelector() squirrel.SelectBuilder {
	return psql.Select(
		"team.team_id",
		"team.name",
		"team.description",
		"team.tag",
		"team.icon_small",
		"team.icon_large",
		"team.wins",
		"team.losses",
		"player.player_id",
		"player.external_id",
		"player.main_roster",
		"player.position",
	).
		From("team").
		LeftJoin("player ON team.team_id = player.team_id")
}

func GetScannedLoLTeamStub(rows *sql.Rows) (*dataModel.LoLTeamStub, error) {
	defer rows.Close()

	var team dataModel.LoLTeamStub

	for rows.Next() {
		var (
			playerId         sql.NullInt64
			playerExternalId sql.NullString
			playerMainRoster sql.NullBool
			playerPosition   sql.NullString
		)
		if err := rows.Scan(
			&team.TeamId,
			&team.Name,
			&team.Description,
			&team.Tag,
			&team.IconSmall,
			&team.IconLarge,
			&team.Wins,
			&team.Losses,
			&playerId,
			&playerExternalId,
			&playerMainRoster,
			&playerPosition,
		); err != nil {
			return nil, err
		}
		if playerExternalId.Valid && playerExternalId.String != "" {
			if playerMainRoster.Bool {
				team.MainRoster = append(team.MainRoster, &dataModel.LoLPlayerStub{
					PlayerId:   int(playerId.Int64),
					ExternalId: playerExternalId.String,
					MainRoster: playerMainRoster.Bool,
					Position:   playerPosition.String,
				})
			} else {
				team.SubstituteRoster = append(team.SubstituteRoster, &dataModel.LoLPlayerStub{
					PlayerId:   int(playerId.Int64),
					ExternalId: playerExternalId.String,
					MainRoster: playerMainRoster.Bool,
					Position:   playerPosition.String,
				})
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &team, nil
}

func getTeamWithManagersSelector() squirrel.SelectBuilder {
	return psql.Select(
		"team.team_id",
		"team.name",
		"team.tag",
		"team.icon_small",
		"user_.user_id",
		"user_.email",
		"team_permissions.administrator",
		"team_permissions.information",
		"team_permissions.games",
	).
		From("team_permissions").
		LeftJoin("user_ ON team_permissions.user_id = user_.user_id").
		LeftJoin("team ON team_permissions.team_id = team.team_id")
}

//
func GetScannedAllTeamWithManagers(rows *sql.Rows) ([]*dataModel.TeamWithManagers, error) {
	teams := make([]*dataModel.TeamWithManagers, 0)
	getUniqueTeam := func(newTeam *dataModel.TeamWithManagers) *dataModel.TeamWithManagers {
		for _, team := range teams {
			if newTeam.TeamId == team.TeamId {
				return team
			}
		}
		teams = append(teams, newTeam)
		return newTeam
	}

	defer rows.Close()
	for rows.Next() {
		var team dataModel.TeamWithManagers
		var manager dataModel.TeamManager
		if err := rows.Scan(
			&team.TeamId,
			&team.Name,
			&team.Tag,
			&team.IconSmall,
			&manager.UserId,
			&manager.Email,
			&manager.Administrator,
			&manager.Information,
			&manager.Games,
		); err != nil {
			return nil, err
		}

		uniqueTeam := getUniqueTeam(&team)
		uniqueTeam.Managers = append(uniqueTeam.Managers, &manager)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *TeamWithPlayersArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamWithPlayers(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

// TeamDisplay
type TeamDisplayArray struct {
	rows []*dataModel.TeamDisplay
}

func getTeamDisplaySelector() squirrel.SelectBuilder {
	return psql.Select(
		"team_id",
		"name",
		"tag",
		"icon_small",
	).From("team")
}

func GetScannedTeamDisplay(rows squirrel.RowScanner) (*dataModel.TeamDisplay, error) {
	var team dataModel.TeamDisplay
	if err := rows.Scan(
		&team.TeamId,
		&team.Name,
		&team.Tag,
		&team.IconSmall,
	); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (r *TeamDisplayArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamDisplay(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

// TeamPermissions
func getTeamPermissionsSelector() squirrel.SelectBuilder {
	return psql.Select(
		"team.team_id",
		"team.name",
		"team.tag",
		"team.icon_small",
		"team_permissions.administrator",
		"team_permissions.information",
		"team_permissions.games",
	).
		From("team").
		Join("team_permissions ON team.team_id = team_permissions.team_id")
}

type TeamPermissionsArray struct {
	rows []*dataModel.TeamPermissions
}

func GetScannedTeamPermissions(rows squirrel.RowScanner) (*dataModel.TeamPermissions, error) {
	var teamPermissions dataModel.TeamPermissions
	if err := rows.Scan(
		&teamPermissions.TeamId,
		&teamPermissions.Name,
		&teamPermissions.Tag,
		&teamPermissions.IconSmall,
		&teamPermissions.Administrator,
		&teamPermissions.Information,
		&teamPermissions.Games,
	); err != nil {
		return nil, err
	} else {
		return &teamPermissions, nil
	}
}

func (r *TeamPermissionsArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamPermissions(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}
