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
		//
		//available := make([]int, len(teams))
		//
		//played := make([][]bool, len(teams))
		//for i := range played {
		//	played[i] = make([]bool, len(teams))
		//}
		//
		//for _, i := range teams {
		//	for _, j := range teams {
		//		played[i][j] = i == j
		//	}
		//}
		//
		//for a := 0; a < len(teams) - 1; a++ {
		//	for _, i := range teams {
		//		available[i] = -1
		//	}
		//
		//	for teamNumber := range teams {
		//		availableCursor := 0
		//		fmt.Printf("%v\n", available)
		//		println(played[0][2])
		//		for availableCursor < len(teams) {
		//			if available[availableCursor] == -1 && math.Mod(float64(availableCursor), 2) == 0 {
		//				available[availableCursor] = teamNumber
		//				break
		//			} else if available[availableCursor] == -1 && math.Mod(float64(availableCursor), 2) == 1 {
		//				if !played[available[availableCursor-1]][teamNumber] {
		//					available[availableCursor] = teamNumber
		//					break
		//				}
		//			}
		//			availableCursor++
		//		}
		//	}
		//
		//	for i := 0; i < len(teams); i += 2 {
		//		requiredGames = append(requiredGames,
		//			sgame{available[i], available[i+1]})
		//		println("plays ", i)
		//		played[available[i]][available[i+1]] = true
		//		played[available[i+1]][available[i]] = true
		//	}
		//}















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
