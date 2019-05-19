package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// In and out of DAO

type LeagueInformationDTO struct {
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

type LeagueInformationArray struct {
	rows []*LeagueInformationDTO
}

func GetScannedLeagueInformationDTO(rows squirrel.RowScanner) (*LeagueInformationDTO, error) {
	var leagueInformation LeagueInformationDTO
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

func (r *LeagueInformationArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedLeagueInformationDTO(rows)
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
