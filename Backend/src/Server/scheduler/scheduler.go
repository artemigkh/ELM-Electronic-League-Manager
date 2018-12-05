package scheduler

import (
	"fmt"
	"sort"
	"time"
)

func (s *Scheduler) InitScheduler(tournamentType int, stretch bool, gameDuration time.Duration, start, end time.Time, teams []int) {
	s.tournamentType = tournamentType
	s.stretch = stretch
	s.gameDuration = gameDuration
	s.start = start
	s.end = end
	s.teams = teams
}

func (s *Scheduler) AddDailyAvailability(hour, minute int, duration time.Duration) {

}

func (s *Scheduler) AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration) {
	weekCursor := time.Time(s.start)
	for weekCursor.Weekday() != dayOfWeek {
		weekCursor = weekCursor.Add(time.Hour * 24)
	}

	for weekCursor.Before(s.end) {
		s.availabilities = append(s.availabilities, Availability{
			weekCursor.Add(time.Hour*time.Duration(hour) +
				time.Minute*time.Duration(minute)),
			weekCursor.Add(time.Hour*time.Duration(hour) +
				time.Minute*time.Duration(minute) +
				duration),
		})
		weekCursor = weekCursor.AddDate(0, 0, 7)
	}
}

func (s *Scheduler) GetSchedule() []Game {
	sort.Slice(s.availabilities, func(i, j int) bool {
		return s.availabilities[i].Start.Before(s.availabilities[j].Start)
	})
	for _, avail := range s.availabilities {
		fmt.Printf("%v, %v\n", avail.Start.Format(time.UnixDate), avail.End.Format(time.UnixDate))
	}

	// if no availabilities, fail
	if len(s.availabilities) == 0 {
		return nil
	}

	scheduleCursor := s.availabilities[0].Start
	availabilityIndex := 0
	var games []Game
	for _, game := range getRequiredGames(s.tournamentType, s.teams) {
		games = append(games, Game{game.team1, game.team2, scheduleCursor.Unix()})
		scheduleCursor = scheduleCursor.Add(s.gameDuration)
		if scheduleCursor.Add(s.gameDuration).After(s.availabilities[availabilityIndex].End) {
			if availabilityIndex+1 >= len(s.availabilities) {
				return nil
			} else {
				availabilityIndex++
				scheduleCursor = s.availabilities[availabilityIndex].Start
			}
		}
	}

	return games
}
