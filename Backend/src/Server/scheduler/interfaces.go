package scheduler

import "time"

const (
	RoundRobin       = iota
	DoubleRoundRobin = iota
)

type Game struct {
	Team1Id  int
	Team2Id  int
	GameTime int
}

var tournamentStringToEnum = map[string]int{
	"roundrobin":       RoundRobin,
	"doubleroundrobin": DoubleRoundRobin,
}

var weekdays = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}

type IScheduler interface {
	GetTournamentFromString(tournament string) int
	GetWeekdayFromString(weekday string) time.Weekday

	InitScheduler(tournamentType, roundsPerWeek, concurrentGameNum int, gameDuration time.Duration, start, end time.Time, teams []int)
	AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration)
	GetSchedule() ([]Game, error)

	IsTournamentTypeSupported(tournament string) bool
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
