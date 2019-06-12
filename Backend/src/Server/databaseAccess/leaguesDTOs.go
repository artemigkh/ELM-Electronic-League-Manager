package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type SchedulingAvailabilityDTO struct {
	Id          int  `json:"id"`
	Weekday     int  `json:"weekday"`
	Timezone    int  `json:"timezone"`
	Hour        int  `json:"hour"`
	Minute      int  `json:"minute"`
	Duration    int  `json:"duration"`
	Constrained bool `json:"constrained"`
	Start       int  `json:"start"`
	End         int  `json:"end"`
}

type SchedulingAvailabilityArray struct {
	rows []*SchedulingAvailabilityDTO
}

func GetScannedSchedulingAvailabilityDTO(rows squirrel.RowScanner) (*SchedulingAvailabilityDTO, error) {
	var schedulingInformation SchedulingAvailabilityDTO
	if err := rows.Scan(
		&schedulingInformation.Id,
		&schedulingInformation.Weekday,
		&schedulingInformation.Timezone,
		&schedulingInformation.Hour,
		&schedulingInformation.Minute,
		&schedulingInformation.Duration,
		&schedulingInformation.Constrained,
		&schedulingInformation.Start,
		&schedulingInformation.End,
	); err != nil {
		return nil, err
	} else {
		return &schedulingInformation, nil
	}
}

func (r *SchedulingAvailabilityArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedSchedulingAvailabilityDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

type ManagerDTO struct {
	User        UserDTO             `json:"user"`
	Permissions TeamPermissionsCore `json:"permissions"`
}

type TeamManagerDTO struct {
	Team     *TeamDTO      `json:"team"`
	Managers []*ManagerDTO `json:"managers"`
}

type TeamManagerDTOArray struct {
	rows []*TeamManagerDTO
}

func (r *TeamManagerDTOArray) Scan(rows *sql.Rows) error {
	var user UserDTO
	var team TeamDTO
	var teamPermissions TeamPermissionsCore

	if err := rows.Scan(
		&user.UserId,
		&user.Email,
		&team.Id,
		&team.Name,
		&team.Tag,
		&team.Description,
		&team.IconSmall,
		&teamPermissions.Administrator,
		&teamPermissions.Information,
		&teamPermissions.Games,
	); err != nil {
		return err
	}

	var teamManagerInstance *TeamManagerDTO
	exists := false
	// get TeamManagerDTO instance to insert into if exists
	for _, teamManager := range r.rows {
		if teamManager.Team.Id == team.Id {
			teamManagerInstance = teamManager
			exists = true
		}
	}

	if !exists {
		teamManagerInstance.Team = &team
		r.rows = append(r.rows, teamManagerInstance)
	}

	teamManagerInstance.Managers = append(teamManagerInstance.Managers,
		&ManagerDTO{
			User:        user,
			Permissions: teamPermissions,
		})

	return nil
}
