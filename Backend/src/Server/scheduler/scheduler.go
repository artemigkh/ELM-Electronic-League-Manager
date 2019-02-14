package scheduler

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func (s *Scheduler) InitScheduler(tournamentType, roundsPerWeek, concurrentGameNum int, gameDuration time.Duration, start, end time.Time, teams []int) {
	s.tournamentType = tournamentType
	s.roundsPerWeek = roundsPerWeek
	s.concurrentGameNum = concurrentGameNum
	s.gameDuration = gameDuration
	s.start = start
	s.end = end
	s.teams = teams
}

func (s *Scheduler) AddDailyAvailability(hour, minute int, duration time.Duration) {

}

func leq(t1, t2 time.Time) bool {
	return t1.Before(t2) || t1.Equal(t2)
}

func (s *Scheduler) AddWeeklyAvailability(dayOfWeek time.Weekday, hour, minute int, duration time.Duration) {
	weekCursor := time.Time(s.start)
	for weekCursor.Weekday() != dayOfWeek {
		weekCursor = weekCursor.Add(time.Hour * 24)
	}

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

func (s *Scheduler) GetSchedule() []Game {
	sort.Slice(s.gameBlocks, func(i, j int) bool {
		return s.gameBlocks[i].Start.Before(s.gameBlocks[j].Start)
	})
	for _, avail := range s.gameBlocks {
		fmt.Printf("%v, %v, %v\n", avail.Start.Format(time.UnixDate), avail.End.Format(time.UnixDate), avail.teams)
	}

	// if no gameBlocks, fail
	if len(s.gameBlocks) == 0 {
		return nil
	}

	// split up game blocks into weeks
	weekGameBlocks := make([][]*GameBlock, 0)
	blocks := make([]*GameBlock, 0)
	weekCursor := time.Time(s.start)
	for _, block := range s.gameBlocks {
		if block.Start.After(weekCursor.AddDate(0, 0, 7)) {
			weekGameBlocks = append(weekGameBlocks, blocks)
			weekCursor = weekCursor.AddDate(0, 0, 7)
			blocks = make([]*GameBlock, 0)
		}
		blocks = append(blocks, block)
	}

	println("week game blocks")
	for weeknum, week := range weekGameBlocks {
		fmt.Printf("week: %v\n", weeknum)
		for _, block := range week {
			fmt.Printf("start: %v\n", block.Start.Format(time.UnixDate))
		}
	}

	requiredGames := getRequiredGames(s.tournamentType, s.teams)
	if len(requiredGames) > len(s.gameBlocks)*s.concurrentGameNum {
		return nil
	}

	var games []Game

	if s.roundsPerWeek <= 0 {
		for _, game := range requiredGames {
			scheduled := false
			for i, block := range s.gameBlocks {
				if in(game.team1, block.teams) && in(game.team2, block.teams) {
					scheduled = true
					games = append(games, Game{game.team1, game.team2, block.Start.Unix()})
					block.NumGames += 1

					if block.NumGames >= s.concurrentGameNum {
						s.gameBlocks = append(s.gameBlocks[:i], s.gameBlocks[i+1:]...)
					}
					break
				}
			}
			if !scheduled {
				return nil
			}
		}
	} else {

		gameIndex := 0

		gamesPerRound := int(math.Ceil(float64(len(s.teams)) / 2))
		gamesPerWeek := gamesPerRound * s.roundsPerWeek
		totalRounds := len(requiredGames) / gamesPerRound
		totalWeeks := totalRounds / s.roundsPerWeek

		if totalWeeks > len(weekGameBlocks) {
			return nil
		}

		for _, weekBlocks := range weekGameBlocks {
			if gamesPerWeek > len(weekBlocks) {
				return nil
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
							block.Start.Unix()})
						block.NumGames += 1

						if block.NumGames >= s.concurrentGameNum {
							weekBlocks = append(weekBlocks[:i], weekBlocks[i+1:]...)
						}
						break
					}
				}
				if !scheduled {
					return nil
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

	return games
}
