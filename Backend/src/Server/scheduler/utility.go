package scheduler

import (
	"math/rand"
	"time"
)

type sgame struct {
	team1 int
	team2 int
}

func getRequiredGames(tournamentType int, teams []int) []sgame {
	var requiredGames []sgame
	numTeams := len(teams)
	if tournamentType == RoundRobin {
		curPos := 0
		curJmp := 1
		for numGames := 0; numGames < (numTeams*(numTeams-1))/2; numGames++ {
			requiredGames = append(requiredGames,
				sgame{teams[curPos], teams[(curPos+curJmp)%numTeams]})
			curPos++
			if curPos == numTeams {
				curPos = 0
				curJmp++
			}
		}

		//shuffle order of games
		r := rand.New(rand.NewSource(time.Now().Unix()))
		temp := make([]sgame, len(requiredGames))
		perm := r.Perm(len(requiredGames))
		for i, randIndex := range perm {
			temp[i] = requiredGames[randIndex]
		}
		requiredGames = temp
	} else if tournamentType == DoubleRoundRobin {
		curPos := 0
		curJmp := 1
		for numGames := 0; numGames < (numTeams*(numTeams-1))/2; numGames++ {
			requiredGames = append(requiredGames,
				sgame{teams[curPos], teams[(curPos+curJmp)%numTeams]})
			curPos++
			if curPos == numTeams {
				curPos = 0
				curJmp++
			}
		}
		var final []sgame
		//shuffle order of games
		r := rand.New(rand.NewSource(time.Now().Unix()))
		temp := make([]sgame, len(requiredGames))
		perm := r.Perm(len(requiredGames))
		for i, randIndex := range perm {
			temp[i] = requiredGames[randIndex]
		}
		final = temp

		temp = make([]sgame, len(requiredGames))
		perm = r.Perm(len(requiredGames))
		for i, randIndex := range perm {
			temp[i] = requiredGames[randIndex]
		}
		final = append(final, temp...)
		requiredGames = final
	}
	return requiredGames
}
