package IntegrationTests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
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