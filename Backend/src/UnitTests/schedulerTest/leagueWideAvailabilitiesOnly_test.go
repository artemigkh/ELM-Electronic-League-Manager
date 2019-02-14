package schedulerTest

import (
	"Server/scheduler"
	"fmt"
	"testing"
	"time"
)

func Test_LeagueWideWeeklyAvailabilities(t *testing.T) {
	s := scheduler.Scheduler{}
	est, _ := time.LoadLocation("America/New_York")
	s.InitScheduler(scheduler.DoubleRoundRobin, 2, 1, time.Hour,
		time.Date(2018, time.October, 15, 0, 0, 0, 0, est),
		time.Date(2018, time.December, 23, 0, 0, 0, 0, est),
		[]int{1, 2, 3, 4, 5, 6, 7, 8})
	s.AddWeeklyAvailability(time.Friday, 12+6, 0, time.Hour*2)
	s.AddWeeklyAvailability(time.Saturday, 12+4, 0, time.Hour*6)
	s.AddWeeklyAvailability(time.Sunday, 12+5, 0, time.Hour*5)
	games := s.GetSchedule()
	for _, game := range games {
		fmt.Printf("%v vs %v - %v\n", game.Team1Id, game.Team2Id, time.Unix(game.GameTime, 0).Format(time.UnixDate))
	}
}
