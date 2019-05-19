package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type LeagueDTO struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Game        string `json:"game"`
	PublicView  bool   `json:"publicView"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signupStart"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
}

type LeagueDTOArray struct {
	rows []*LeagueDTO
}

func GetScannedLeagueDTO(rows squirrel.RowScanner) (*LeagueDTO, error) {
	var leagueInformation LeagueDTO
	if err := rows.Scan(
		&leagueInformation.Id,
		&leagueInformation.Name,
		&leagueInformation.Description,
		&leagueInformation.Game,
		&leagueInformation.PublicView,
		&leagueInformation.PublicJoin,
		&leagueInformation.SignupStart,
		&leagueInformation.SignupEnd,
		&leagueInformation.LeagueStart,
		&leagueInformation.LeagueEnd,
	); err != nil {
		return nil, err
	} else {
		return &leagueInformation, nil
	}
}

func (r *LeagueDTOArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedLeagueDTO(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

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

type LeaguePermissionsDTO struct {
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

func GetScannedLeaguePermissionsDTO(rows squirrel.RowScanner) (*LeaguePermissionsDTO, error) {
	var leaguePermissions LeaguePermissionsDTO
	if err := rows.Scan(
		&leaguePermissions.Administrator,
		&leaguePermissions.CreateTeams,
		&leaguePermissions.EditTeams,
		leaguePermissions.EditGames,
	); err != nil {
		return nil, err
	} else {
		return &leaguePermissions, nil
	}
}

type ManagerDTO struct {
	User        UserDTO            `json:"user"`
	Permissions TeamPermissionsDTO `json:"permissions"`
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
	var teamPermissions TeamPermissionsDTO

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
		&teamPermissions.Players,
		&teamPermissions.ReportResults,
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
