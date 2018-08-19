package IntegrationTests

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/Pallinder/go-randomdata"
)

func checkGameAgainstRepresentation(t *testing.T, g *game) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET",
		fmt.Sprintf("api/games/%v", g.Id), 200)

	assert.Equal(t, responseMap["id"], g.Id)
	assert.Equal(t, responseMap["leagueId"], g.LeagueId)
	assert.Equal(t, responseMap["team1Id"], g.Team1Id)
	assert.Equal(t, responseMap["team2Id"], g.Team2Id)
	assert.Equal(t, responseMap["gameTime"], g.GameTime)
	assert.Equal(t, responseMap["complete"], g.Complete)
	assert.Equal(t, responseMap["winnerId"], g.WinnerId)
	assert.Equal(t, responseMap["scoreTeam1"], g.ScoreTeam1)
	assert.Equal(t, responseMap["scoreTeam2"], g.ScoreTeam2)
}

func reportGame(t *testing.T, g *game, winnerId, scoreTeam1, scoreTeam2 float64) {
	body := make(map[string]interface{})
	body["winnerId"] = winnerId
	body["scoreTeam1"] = scoreTeam1
	body["scoreTeam2"] = scoreTeam2

	makeApiCall(t, body, "POST", fmt.Sprintf("api/games/report/%v", g.Id), 200)
}

func randomlyDecideAndReportGame(t *testing.T, g *game) {
	scoreTeam1 := float64(randomdata.Number(0, 3))
	scoreTeam2 := float64(randomdata.Number(0, 3))

	for scoreTeam1 == scoreTeam2 {
		scoreTeam2 = float64(randomdata.Number(0, 3))
	}

	var winnerId float64

	if scoreTeam1 > scoreTeam2 {
		winnerId = g.Team1Id
	} else {
		winnerId = g.Team2Id
	}

	g.Complete = true
	g.WinnerId = winnerId
	g.ScoreTeam1 = scoreTeam1
	g.ScoreTeam2 = scoreTeam2
	reportGame(t, g, winnerId, scoreTeam1, scoreTeam2)
}