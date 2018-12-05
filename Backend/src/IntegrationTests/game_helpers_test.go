package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"testing"
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

func randomlyDecideAndReportGame(t *testing.T, g *game, teams []*team) {
	var t1Strength int
	var t2Strength int
	for _, team := range teams {
		if team.Id == g.Team1Id {
			t1Strength = team.Strength
		} else if team.Id == g.Team2Id {
			t2Strength = team.Strength
		}
	}
	totStrength := t1Strength + t2Strength
	t1Score := 0
	t2Score := 0

	for t1Score < 3 && t2Score < 3 {
		if randomdata.Number(totStrength) < t1Strength {
			t1Score++
		} else {
			t2Score++
		}
	}

	var winnerId float64

	if t1Score > t2Score {
		winnerId = g.Team1Id
	} else {
		winnerId = g.Team2Id
	}

	g.Complete = true
	g.WinnerId = winnerId
	g.ScoreTeam1 = float64(t1Score)
	g.ScoreTeam2 = float64(t2Score)
	reportGame(t, g, winnerId, g.ScoreTeam1, g.ScoreTeam2)
}
