package scheduler

import "time"

const (
	RoundRobin       = iota
	DoubleRoundRobin = iota
)

type Game struct {
	Team1Id  int
	Team2Id  int
	GameTime int64
}

type IScheduler interface {
	InitScheduler(tournamentType, roundsPerWeek, concurrentGameNum int, gameDuration time.Duration, start, end time.Time, teams []int)
	AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration)
	GetSchedule() []Game
}

type GameBlock struct {
	Start    time.Time
	End      time.Time
	NumGames int
	teams    []int
}

type Scheduler struct {
	tournamentType    int
	roundsPerWeek     int
	concurrentGameNum int
	gameDuration      time.Duration
	start             time.Time
	end               time.Time
	teams             []int
	gameBlocks        []*GameBlock
}
