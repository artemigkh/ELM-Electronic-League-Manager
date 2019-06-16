package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

//Availability
func (avail *AvailabilityCore) Validate(leagueId int) (bool, string, error) {
	//TODO: implement these validate functions
	return validate(avail.timestamps())
}

func (avail *AvailabilityCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
}

type AvailabilityArray struct {
	rows []*Availability
}

func getAvailabilitySelector(leagueId int) squirrel.SelectBuilder {
	return psql.Select(
		"availability_id",
		"start_time",
		"end_time",
	).
		From("availability").
		Where("is_recurring_weekly = false AND league_id = ?", leagueId)

}

func GetScannedAvailability(rows squirrel.RowScanner) (*Availability, error) {
	var availability Availability
	if err := rows.Scan(
		&availability.AvailabilityId,
		&availability.StartTime,
		&availability.EndTime,
	); err != nil {
		return nil, err
	} else {
		return &availability, nil
	}
}

func (r *AvailabilityArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedAvailability(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

//WeeklyAvailability
func (avail *WeeklyAvailabilityCore) validate(leagueId, availabilityId int) (bool, string, error) {
	//TODO: implement these validate functions
	return validate(avail.timestamps())
}

func (avail *WeeklyAvailabilityCore) ValidateNew(leagueId int) (bool, string, error) {
	return avail.validate(leagueId, 0)
}

func (avail *WeeklyAvailabilityCore) ValidateEdit(leagueId, availabilityId int) (bool, string, error) {
	return avail.validate(leagueId, availabilityId)
}

func (avail *WeeklyAvailabilityCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
}

type WeeklyAvailabilityArray struct {
	rows []*WeeklyAvailability
}

func getWeeklyAvailabilitySelector(leagueId int) squirrel.SelectBuilder {
	return psql.Select(
		"availability.availability_id",
		"availability.start_time",
		"availability.end_time",
		"weekly_recurrence.weekday",
		"weekly_recurrence.timezone",
		"weekly_recurrence.hour",
		"weekly_recurrence.minute",
		"weekly_recurrence.duration",
	).
		From("availability").
		Join("weekly_recurrence ON availability.availability_id = weekly_recurrence.availability_id").
		Where("is_recurring_weekly = true AND league_id = ?", leagueId)
}

func GetScannedWeeklyAvailability(rows squirrel.RowScanner) (*WeeklyAvailability, error) {
	var availability WeeklyAvailability
	if err := rows.Scan(
		&availability.AvailabilityId,
		&availability.StartTime,
		&availability.EndTime,
		&availability.Weekday,
		&availability.Timezone,
		&availability.Hour,
		&availability.Minute,
		&availability.Duration,
	); err != nil {
		return nil, err
	} else {
		return &availability, nil
	}
}

func (r *WeeklyAvailabilityArray) Scan(rows *sql.Rows) error {
	row, err := GetScannedWeeklyAvailability(rows)
	if err != nil {
		return err
	} else {
		r.rows = append(r.rows, row)
		return nil
	}
}

// SchedulingParameters
func (params *SchedulingParameters) Validate() (bool, string, error) {
	//TODO: implement these validate functions
	return validate(params.tournamentType())
}

func (params *SchedulingParameters) tournamentType() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
}
