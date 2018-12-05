package scheduler

import "time"

const (
	RoundRobin       = iota
	DoubleRoundRobin = iota
	SingleElim       = iota
	DoubleElim       = iota
)

type Game struct {
	Team1Id  int
	Team2Id  int
	GameTime int64
}

type IScheduler interface {
	InitScheduler(tournamentType int, stretch bool, gameDuration time.Duration, start, end time.Time, teams []int)
	AddDailyAvailability(hour, minute int, duration time.Duration)
	AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration)
	GetSchedule() []Game
}

type Availability struct {
	Start time.Time
	End   time.Time
}

type Scheduler struct {
	tournamentType int
	stretch        bool
	gameDuration   time.Duration
	start          time.Time
	end            time.Time
	teams          []int
	availabilities []Availability
}
