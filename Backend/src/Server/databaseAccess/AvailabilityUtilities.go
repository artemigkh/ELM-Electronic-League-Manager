package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type AvailabilityArray struct {
	rows []*dataModel.Availability
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

func GetScannedAvailability(rows squirrel.RowScanner) (*dataModel.Availability, error) {
	var availability dataModel.Availability
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

type WeeklyAvailabilityArray struct {
	rows []*dataModel.WeeklyAvailability
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

func GetScannedWeeklyAvailability(rows squirrel.RowScanner) (*dataModel.WeeklyAvailability, error) {
	var availability dataModel.WeeklyAvailability
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
