package scheduler

type sgame struct {
	team1 int
	team2 int
}

func roundRobin(teams []int) []sgame {
	var requiredGames []sgame
	startIndex := 0
	if len(teams)%2 == 1 {
		// Set bye to position 0 and don't use it
		teams = append([]int{-1}, teams...)
		startIndex = 1
	}
	for round := 0; round < len(teams)-1; round++ {
		for i := startIndex; i < len(teams)/2; i++ {
			requiredGames = append(requiredGames,
				sgame{teams[i], teams[len(teams)-1-i]})
		}
		teams = append([]int{teams[0], teams[len(teams)-1]}, teams[1:len(teams)-1]...)
	}
	return requiredGames
}

func getRequiredGames(tournamentType int, teams []int) []sgame {
	if tournamentType == RoundRobin {
		return roundRobin(teams)
	} else if tournamentType == DoubleRoundRobin {
		return append(roundRobin(teams), roundRobin(teams)...)
	}
	return nil
}
