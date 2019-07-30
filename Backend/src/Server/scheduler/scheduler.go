package scheduler

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"sort"
	"strings"
	"time"
)

func (s *Scheduler) GetTournamentFromString(tournament string) int {
	return tournamentStringToEnum[tournament]
}
func (s *Scheduler) IsTournamentTypeSupported(tournament string) bool {
	_, ok := tournamentStringToEnum[strings.ToLower(tournament)]
	return ok
}

func (s *Scheduler) GetWeekdayFromString(weekday string) time.Weekday {
	return weekdays[weekday]
}

func (s *Scheduler) InitScheduler(tournamentType, roundsPerWeek, concurrentGameNum int,
	gameDuration time.Duration, start, end time.Time, teams []int) {

	s.tournamentType = tournamentType
	s.roundsPerWeek = roundsPerWeek
	s.concurrentGameNum = concurrentGameNum
	s.gameDuration = gameDuration
	s.start = start
	s.end = end
	s.teams = teams
}

func leq(t1, t2 time.Time) bool {
	return t1.Before(t2) || t1.Equal(t2)
}

func (s *Scheduler) AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration) {
	fmt.Printf("League start: %v\n", s.start.Format(time.UnixDate))
	weekCursor := time.Time(s.start)
	weekCursor = time.Date(
		weekCursor.Year(),
		weekCursor.Month(),
		weekCursor.Day(),
		hour,
		minute,
		0,
		0,
		weekCursor.Location(),
	)
	for weekCursor.Weekday() != dayOfWeek {
		weekCursor = weekCursor.Add(time.Hour * 24)
	}
	fmt.Printf("week cursor start: %v\n", s.start.Format(time.UnixDate))
	fmt.Printf("League end: %v\n", s.end.Format(time.UnixDate))
	for weekCursor.Before(s.end) {
		blockCursor := weekCursor.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))
		for leq(blockCursor.Add(s.gameDuration), weekCursor.Add(time.Hour*time.Duration(hour)+
			time.Minute*time.Duration(minute)+duration)) {

			s.gameBlocks = append(s.gameBlocks, &GameBlock{
				blockCursor,
				blockCursor.Add(s.gameDuration), 0,
				s.teams,
			})
			blockCursor = blockCursor.Add(s.gameDuration)
		}

		weekCursor = weekCursor.AddDate(0, 0, 7)
		fmt.Printf("week cursor current: %v\n", weekCursor.Format(time.UnixDate))
		fmt.Printf("League end: %v\n", s.end.Format(time.UnixDate))
	}
}

func in(el int, list []int) bool {
	for _, a := range list {
		if a == el {
			return true
		}
	}
	return false
}

func (s *Scheduler) GetSchedule() ([]Game, error) {
	sort.Slice(s.gameBlocks, func(i, j int) bool {
		return s.gameBlocks[i].Start.Before(s.gameBlocks[j].Start)
	})
	for _, avail := range s.gameBlocks {
		fmt.Printf("%v, %v, %v\n", avail.Start.Format(time.UnixDate), avail.End.Format(time.UnixDate), avail.teams)
	}

	// if no gameBlocks, fail
	if len(s.gameBlocks) == 0 {
		return nil, errors.New("Zero Available Game Blocks")
	}

	// split up game blocks into weeks
	weekGameBlocks := make([][]*GameBlock, 0)
	blocks := make([]*GameBlock, 0)
	weekCursor := time.Time(s.start)
	for _, block := range s.gameBlocks {
		if block.Start.After(weekCursor.AddDate(0, 0, 7)) {
			println(fmt.Sprintf("Updating week cursor to %v\n", weekCursor.AddDate(0, 0, 7).Format(time.UnixDate)))
			weekGameBlocks = append(weekGameBlocks, blocks)
			weekCursor = weekCursor.AddDate(0, 0, 7)
			blocks = make([]*GameBlock, 0)
		}
		blocks = append(blocks, block)
	}
	weekGameBlocks = append(weekGameBlocks, blocks)

	println("week game blocks")
	for weeknum, week := range weekGameBlocks {
		fmt.Printf("week: %v\n", weeknum)
		for _, block := range week {
			fmt.Printf("start: %v\n", block.Start.Format(time.UnixDate))
		}
	}

	requiredGames := getRequiredGames(s.tournamentType, s.teams)
	if len(requiredGames) > len(s.gameBlocks)*s.concurrentGameNum {
		return nil, errors.New(fmt.Sprintf("Number of available game blocks (%v) is "+
			"less than number of required games (%v).", len(s.gameBlocks)*s.concurrentGameNum, len(requiredGames)))
	}

	var games []Game

	if s.roundsPerWeek <= 0 {
		for _, game := range requiredGames {
			scheduled := false
			for i, block := range s.gameBlocks {
				if in(game.team1, block.teams) && in(game.team2, block.teams) {
					scheduled = true
					games = append(games, Game{game.team1, game.team2, int(block.Start.Unix())})
					block.NumGames += 1

					if block.NumGames >= s.concurrentGameNum {
						s.gameBlocks = append(s.gameBlocks[:i], s.gameBlocks[i+1:]...)
					}
					break
				}
			}
			if !scheduled {
				return nil, errors.New("Scheduling failed due to constraints")
			}
		}
	} else {

		gameIndex := 0

		gamesPerRound := int(math.Ceil(float64(len(s.teams)) / 2))
		gamesPerWeek := gamesPerRound * s.roundsPerWeek
		totalRounds := len(requiredGames) / gamesPerRound
		totalWeeks := totalRounds / s.roundsPerWeek

		if totalWeeks > len(weekGameBlocks) {
			return nil, errors.New(fmt.Sprintf("Number of required weeks(%v) "+
				"for specified number of rounds per week is larger than amount of available weeks(%v)",
				totalWeeks, len(weekGameBlocks)))
		}

		for weekNum, weekBlocks := range weekGameBlocks {
			if gamesPerWeek > len(weekBlocks) {
				return nil, errors.New(fmt.Sprintf("Number of game blocks(%v) on week %v is"+
					"smaller than number of required games for week %v (%v)",
					len(weekBlocks), weekNum+1, weekNum+1, gamesPerWeek))
			}

			for i := 0; i < gamesPerWeek; i++ {
				// schedule gamesPerWeek sequential games in a row
				gameToSchedule := requiredGames[gameIndex]
				scheduled := false
				for i, block := range weekBlocks {
					if in(gameToSchedule.team1, block.teams) && in(gameToSchedule.team2, block.teams) {
						scheduled = true
						gameIndex++
						games = append(games, Game{
							gameToSchedule.team1,
							gameToSchedule.team2,
							int(block.Start.Unix())})
						block.NumGames += 1

						if block.NumGames >= s.concurrentGameNum {
							weekBlocks = append(weekBlocks[:i], weekBlocks[i+1:]...)
						}
						break
					}
				}
				if !scheduled {
					return nil, errors.New(fmt.Sprintf("Scheduling failed on week %v due to constraints", weekNum+1))
				}
				if gameIndex >= len(requiredGames) {
					break
				}

				//games = append(games, Game{
				//	requiredGames[gameIndex].team1,
				//	requiredGames[gameIndex].team2,
				//	weekBlocks[blockIndex].Start.Unix()})
				//gameIndex++
				//if gameIndex >= len(requiredGames) {
				//	break
				//}
				//
				//blockIndex++
			}
			if gameIndex >= len(requiredGames) {
				break
			}
		}
	}

	return games, nil
}
