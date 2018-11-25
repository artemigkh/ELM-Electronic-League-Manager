package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func checkCantMakeLeagueLoggedOut(t *testing.T) {
	responseMap := makeApiCallAndGetMap(t, nil, "POST", "api/leagues/", 403)
	assert.Equal(t, responseMap["error"].(string), "notLoggedIn")
}

func setActiveLeague(t *testing.T, l *league) {
	//makeApiCall(t, nil, "POST", "api/leagues/setActiveLeague/" + string(l.Id), 200)
	makeApiCall(t, nil, "POST",
		fmt.Sprintf("api/leagues/setActiveLeague/%v", l.Id), 200)
}

func checkCantGetLeagueNoActiveLeague(t *testing.T) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET", "api/leagues/", 403)
	assert.Equal(t, responseMap["error"].(string), "noActiveLeague")
}

func checkLeagueSelected(t *testing.T, l *league) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET", "api/leagues/", 200)
	assert.Equal(t, responseMap["id"], l.Id)
}

func joinLeague(t *testing.T) {
	makeApiCall(t, nil, "POST", "api/leagues/join", 200)
}

func checkTeamsAgainstLeagueSummary(t *testing.T, teams []*team) {
	responseMapArray := makeApiCallAndGetMapArray(t, nil, "GET",
		"api/leagues/teamSummary", 200)

	matchingTeams := 0
	for _, teamSummary := range responseMapArray {
		for _, m := range teams {
			if m.Id == teamSummary["id"] {
				assert.Equal(t, m.Name, teamSummary["name"])
				assert.Equal(t, m.Tag, teamSummary["tag"])
				assert.Equal(t, m.Wins, teamSummary["wins"])
				assert.Equal(t, m.Losses, teamSummary["losses"])

				matchingTeams++
			}
		}
	}
	assert.Equal(t, matchingTeams, len(responseMapArray))
}

func checkGamesAgainstLeagueSummary(t *testing.T, games []*game) {
	responseMapArray := makeApiCallAndGetMapArray(t, nil, "GET",
		"api/leagues/gameSummary", 200)

	for _, element := range responseMapArray {
		for _, game := range games {
			if game.Team1Id == element["id"] {
				assert.Equal(t, game.Team1Id, element["team1Id"])
				assert.Equal(t, game.Team2Id, element["team2Id"])
				assert.Equal(t, game.GameTime, element["gameTime"])
				assert.Equal(t, game.Complete, element["complete"])
				assert.Equal(t, game.WinnerId, element["winnerId"])
				assert.Equal(t, game.ScoreTeam1, element["scoreTeam1"])
				assert.Equal(t, game.ScoreTeam2, element["scoreTeam2"])
			}
		}
	}
}

func randomlyUnscheduleGames(t *testing.T, l *league, n int) {
	for i := 0; i < n; i++ {
		removedIndex := randomdata.Number(0, len(l.Games)-1)
		makeApiCall(t, nil, "DELETE",
			fmt.Sprintf("api/games/%v", l.Games[removedIndex].Id), 200)
		l.Games = append(l.Games[:removedIndex], l.Games[removedIndex+1:]...)
	}
}

func checkTeamStandingsSortedProperly(t *testing.T) {
	responseMapArray := makeApiCallAndGetMapArray(t, nil, "GET",
		"api/leagues/teamSummary", 200)

	previousWins := math.MaxFloat64
	previousLosses := float64(math.MinInt32)
	for _, teamSummary := range responseMapArray {
		assert.True(t, previousWins >= teamSummary["wins"].(float64))
		if previousWins == teamSummary["wins"].(float64) {
			assert.True(t, previousLosses <= teamSummary["losses"].(float64))
		}
		previousWins = teamSummary["wins"].(float64)
		previousLosses = teamSummary["losses"].(float64)
	}
}

func checkLeagueManagersCorrect(t *testing.T, l *league) {
	responseMap := makeApiCallAndGetMapArray(t, nil, "GET", "api/leagues/teamManagers", 200)
	matchingTeams := 0
	for _, team := range responseMap {
		for _, teamRep := range l.Teams {
			if teamRep.Id == team["teamId"] {
				assert.Equal(t, teamRep.Name, team["teamName"])
				assert.Equal(t, teamRep.Tag, team["teamTag"])

				matchingManagers := 0
				for _, manager := range team["managers"].([]interface{}) {
					for _, managerRep := range teamRep.Managers {
						if manager.(map[string]interface{})["userEmail"] == managerRep.Email {
							assert.True(t, manager.(map[string]interface{})["administrator"].(bool))
							assert.True(t, manager.(map[string]interface{})["information"].(bool))
							assert.True(t, manager.(map[string]interface{})["players"].(bool))
							assert.True(t, manager.(map[string]interface{})["reportResults"].(bool))
							matchingManagers++
						}
					}
				}
				assert.Equal(t, matchingManagers, len(teamRep.Managers))
				matchingTeams++
			}
		}
	}
	assert.Equal(t, matchingTeams, len(l.Teams))
}
