package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type TeamDTO struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Tag         string       `json:"tag"`
	Description string       `json:"description"`
	Wins        int          `json:"wins"`
	Losses      int          `json:"losses"`
	IconSmall   string       `json:"iconSmall"`
	IconLarge   string       `json:"iconLarge"`
	Players     []*PlayerDTO `json:"players"`
}

type TeamDTOArray struct {
	rows []*TeamDTO
}

func GetScannedTeamDTO(row squirrel.RowScanner) (*TeamDTO, error) {
	var team TeamDTO
	if err := row.Scan(
		&team.Id,
		&team.Name,
		&team.Tag,
		&team.Description,
		&team.Wins,
		&team.Losses,
		&team.IconSmall,
		&team.IconLarge,
	); err != nil {
		return nil, err
	} else {
		return &team, nil
	}
}

func (r *TeamDTOArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

type TeamPermissionsDTO struct {
	Administrator bool
	Information   bool
	Players       bool
	ReportResults bool
}

type TeamPermissionsDTOArray struct {
	rows []*TeamPermissionsDTO
}

func GetScannedTeamPermissionsDTO(rows squirrel.RowScanner) (*TeamPermissionsDTO, error) {
	var teamPermissions TeamPermissionsDTO
	if err := rows.Scan(
		&teamPermissions.Administrator,
		&teamPermissions.Information,
		&teamPermissions.Players,
		&teamPermissions.ReportResults,
	); err != nil {
		return nil, err
	} else {
		return &teamPermissions, nil
	}
}

func (r *TeamPermissionsDTOArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedTeamPermissionsDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

type PlayerDTO struct {
	Id             int    `json:"id"`
	TeamId         int    `json:"teamId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	ExternalId     string `json:"external_id"`
	Position       string `json:"position"`
	MainRoster     bool   `json:"mainRoster"`
}

type PlayerDTOArray struct {
	rows []*PlayerDTO
}

func GetScannedPlayerDTO(row squirrel.RowScanner) (*PlayerDTO, error) {
	var player PlayerDTO
	if err := row.Scan(
		&player.Id,
		&player.TeamId,
		&player.Name,
		&player.GameIdentifier,
		&player.ExternalId,
		&player.Position,
		&player.MainRoster,
	); err != nil {
		return nil, err
	} else {
		return &player, nil
	}
}

func (r *PlayerDTOArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedPlayerDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}
