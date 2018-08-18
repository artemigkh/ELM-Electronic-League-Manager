package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
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

	tag := randomdata.Letters(4)
	for i := 0; i < len(teams); i++ {
		if tag == teams[i].Tag {
			tag = randomdata.Letters(4)
			i = 0
		}
	}

	body := make(map[string]interface{})
	body["name"] = teamName
	body["tag"] = tag

	return &team {
		Id: makeApiCallAndGetId(t, body, "POST", "api/teams", 200),
		LeagueId: l.Id,
		Name: teamName,
		Tag: tag,
		Wins: 0,
		Losses: 0,
	}
}
