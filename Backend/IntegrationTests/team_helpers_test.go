package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkTeamCreated(t *testing.T, m *team) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET",
		fmt.Sprintf("api/teams/%v", m.Id), 200)

	assert.Equal(t, responseMap["name"].(string), m.Name)
	assert.Equal(t, responseMap["tag"].(string), m.Tag)
	assert.Equal(t, responseMap["wins"], float64(0))
	assert.Equal(t, responseMap["losses"], float64(0))
	assert.Equal(t, responseMap["players"], nil)
}

func addPlayerToTeam(t *testing.T, m *team, l *league, mainRoster bool) {
	playerName := randomdata.FullName(randomdata.RandomGender)

	gameIdentifier := randomdata.SillyName()
	for i := 0; i < len(l.Players); i++ {
		if gameIdentifier == l.Players[i].GameIdentifier {
			gameIdentifier = randomdata.SillyName()
			i = 0
		}
	}

	body := make(map[string]interface{})
	body["teamId"] = m.Id
	body["name"] = playerName
	body["gameIdentifier"] = gameIdentifier
	body["mainRoster"] = mainRoster

	p := &player{
		Id:             makeApiCallAndGetId(t, body, "POST", "api/teams/addPlayer", 200),
		TeamId:         m.Id,
		GameIdentifier: gameIdentifier,
		Name:           playerName,
		mainRoster:     mainRoster,
	}

	m.Players = append(m.Players, p)
	l.Players = append(l.Players, p)
}

func checkTeamAgainstRepresentation(t *testing.T, m *team) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET",
		fmt.Sprintf("api/teams/%v", m.Id), 200)

	assert.Equal(t, responseMap["name"].(string), m.Name)
	assert.Equal(t, responseMap["tag"].(string), m.Tag)
	assert.Equal(t, responseMap["wins"], m.Wins)
	assert.Equal(t, responseMap["losses"], m.Losses)

	matchingPlayers := 0
	for _, playerSummary := range responseMap["players"].([]interface{}) {
		for _, p := range m.Players {
			if p.Id == playerSummary.(map[string]interface{})["id"] {
				assert.Equal(t, p.Name, playerSummary.(map[string]interface{})["name"])
				assert.Equal(t, p.GameIdentifier, playerSummary.(map[string]interface{})["gameIdentifier"])
				assert.Equal(t, p.mainRoster, playerSummary.(map[string]interface{})["mainRoster"])

				matchingPlayers++
			}
		}
	}
	assert.Equal(t, matchingPlayers, len(m.Players))
}
