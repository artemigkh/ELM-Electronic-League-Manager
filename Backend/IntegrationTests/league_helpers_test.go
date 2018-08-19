package IntegrationTests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"math"
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

func checkTeamStandingsSortedProperly(t *testing.T) {
	responseMapArray := makeApiCallAndGetMapArray(t, nil, "GET",
		"api/leagues/teamSummary", 200)

	previousWins := math.MaxFloat64
	previousLosses := float64(math.MinInt32)
	for _, teamSummary := range responseMapArray {
		assert.True(t, previousWins >= teamSummary["wins"].(float64))
		assert.True(t, previousLosses <= teamSummary["losses"].(float64))
		previousWins = teamSummary["wins"].(float64)
		previousLosses = teamSummary["losses"].(float64)
	}
}