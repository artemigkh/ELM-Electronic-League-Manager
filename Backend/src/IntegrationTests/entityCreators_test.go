package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"strings"
	"testing"
)

func createUser(t *testing.T) *user {
	email := randomdata.Email()
	password := randomdata.RandStringRunes(10)

	body := make(map[string]interface{})
	body["email"] = email
	body["password"] = password

	fmt.Printf("creating user with email: %s and password: %s\n", email, password)

	makeApiCall(t, body, "POST", "api/users", 200)
	return &user{
		Email:    email,
		Password: password,
	}
}

func createLeague(t *testing.T, publicView, publicJoin bool) *league {
	leagueName := randomdata.SillyName()

	body := make(map[string]interface{})
	body["name"] = leagueName
	body["publicView"] = publicView
	body["publicJoin"] = publicJoin

	println("creating league with name " + leagueName)

	return &league{
		Id:         makeApiCallAndGetId(t, body, "POST", "api/leagues", 200),
		Name:       leagueName,
		PublicView: publicView,
		PublicJoin: publicJoin,
	}
}

func createTeam(t *testing.T, teams []*team, l *league) *team {
	teamName := randomdata.SillyName()
	for i := 0; i < len(teams); i++ {
		if teamName == teams[i].Name {
			teamName = randomdata.SillyName()
			i = 0
		}
	}

	var tag string
	if randomdata.Boolean() {
		tag = teamName[0:4]
	} else {
		tag = teamName[0:3]
	}

	endNum := 1
	for true {
		unique := true
		for i := 0; i < len(teams); i++ {
			if strings.ToUpper(tag) == teams[i].Tag {
				unique = false
			}
		}
		if unique {
			break
		} else {
			tag = fmt.Sprintf("%v%v", teamName[0:3], endNum)
			endNum++
		}
	}

	tag = strings.ToUpper(tag)

	body := make(map[string]interface{})
	body["name"] = teamName
	body["tag"] = tag

	return &team{
		Id:       makeApiCallAndGetId(t, body, "POST", "api/teams/", 200),
		LeagueId: l.Id,
		Name:     teamName,
		Tag:      tag,
		Wins:     0,
		Losses:   0,
	}
}

func createGame(t *testing.T, l *league, gameTimeInt int, team1Id, team2Id float64) *game {
	gameTime := float64(gameTimeInt)

	body := make(map[string]interface{})
	body["team1Id"] = team1Id
	body["team2Id"] = team2Id
	body["gameTime"] = gameTime

	return &game{
		Id:         makeApiCallAndGetId(t, body, "POST", "api/games", 200),
		LeagueId:   l.Id,
		Team1Id:    team1Id,
		Team2Id:    team2Id,
		GameTime:   gameTime,
		Complete:   false,
		WinnerId:   -1,
		ScoreTeam1: 0,
		ScoreTeam2: 0,
	}
}
