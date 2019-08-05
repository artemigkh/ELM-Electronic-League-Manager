package dataModel

import (
	"Server/scheduler"
	"math"
)

type AvailabilityCore struct {
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`
}

type Availability struct {
	AvailabilityId int `json:"availabilityId"`
	StartTime      int `json:"startTime"`
	EndTime        int `json:"endTime"`
}

func (avail *AvailabilityCore) Validate(leagueId int, leagueDao LeagueDAO) (bool, string, error) {
	return validate(
		validateAvailabilityTimestamps(avail.StartTime, avail.EndTime),
		validateDuringLeague(leagueId, avail.StartTime, leagueDao, AvailabilityStartNotDuringLeague),
		validateDuringLeague(leagueId, avail.EndTime, leagueDao, AvailabilityEndNotDuringLeague))
}

type WeeklyAvailabilityCore struct {
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
	Weekday   string `json:"weekday"`
	Timezone  int    `json:"timezone"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	Duration  int    `json:"duration"`
}

func (avail *WeeklyAvailabilityCore) validate(leagueId, availabilityId int, leagueDao LeagueDAO) (bool, string, error) {
	return validate(
		validateAvailabilityTimestamps(avail.StartTime, avail.EndTime),
		validateDuringLeague(leagueId, avail.StartTime, leagueDao, AvailabilityStartNotDuringLeague),
		validateDuringLeague(leagueId, avail.EndTime, leagueDao, AvailabilityEndNotDuringLeague),
		avail.weekday(),
		avail.timezone(),
		avail.hour(),
		avail.minute(),
		avail.duration())
}

func (avail *WeeklyAvailabilityCore) ValidateNew(leagueId int, leagueDao LeagueDAO) (bool, string, error) {
	return avail.validate(leagueId, 0, leagueDao)
}

func (avail *WeeklyAvailabilityCore) ValidateEdit(leagueId, availabilityId int, leagueDao LeagueDAO) (bool, string, error) {
	return avail.validate(leagueId, availabilityId, leagueDao)
}

func (avail *WeeklyAvailabilityCore) weekday() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		for _, weekday := range [7]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"} {
			if avail.Weekday == weekday {
				return true
			}
		}
		*problemDest = InvalidWeekday
		return false
	}
}

func (avail *WeeklyAvailabilityCore) timezone() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if math.Abs(float64(avail.Timezone)) < 24*3600 {
			return true
		} else {
			*problemDest = InvalidTimezone
			return false
		}
	}
}

func (avail *WeeklyAvailabilityCore) hour() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if avail.Hour < 24 {
			return true
		} else {
			*problemDest = InvalidHour
			return false
		}
	}
}

func (avail *WeeklyAvailabilityCore) minute() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if avail.Minute < 60 {
			return true
		} else {
			*problemDest = InvalidMinute
			return false
		}
	}
}

func (avail *WeeklyAvailabilityCore) duration() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if avail.Duration < 7*24*60 {
			return true
		} else {
			*problemDest = InvalidWeeklyAvailabilityDuration
			return false
		}
	}
}

type WeeklyAvailability struct {
	AvailabilityId int    `json:"availabilityId"`
	StartTime      int    `json:"startTime"`
	EndTime        int    `json:"endTime"`
	Weekday        string `json:"weekday"`
	Timezone       int    `json:"timezone"`
	Hour           int    `json:"hour"`
	Minute         int    `json:"minute"`
	Duration       int    `json:"duration"`
}

type SchedulingParameters struct {
	TournamentType    string `json:"tournamentType"`
	RoundsPerWeek     int    `json:"roundsPerWeek"`
	ConcurrentGameNum int    `json:"concurrentGameNum"`
	GameDuration      int    `json:"gameDuration"`
}

func (params *SchedulingParameters) Validate() (bool, string, error) {
	return validate(params.tournamentType())
}

func (params *SchedulingParameters) tournamentType() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if supported := (&scheduler.Scheduler{}).IsTournamentTypeSupported(params.TournamentType); supported {
			return true
		} else {
			*problemDest = TournamentTypeNotSupported
			return false
		}
	}
}
