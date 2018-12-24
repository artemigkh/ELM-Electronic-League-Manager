package IntegrationTests

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"strings"
	"testing"
	"time"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func createLeague(t *testing.T, publicView, publicJoin bool) *league {
	leagueName := randomdata.SillyName()

	est, _ := time.LoadLocation("America/New_York")
	currentTime := time.Now()
	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, est)
	signupStart := currentTime.AddDate(0, 0, -7*5)
	signupEnd := signupStart.AddDate(0, 0, 7)
	leagueStart := signupEnd.AddDate(0, 0, 7)
	leagueEnd := leagueStart.AddDate(0, 0, 7*7)

	description := ""
	randParagraph := randomdata.Paragraph()
	for len(description+randParagraph) < 500 {
		description += " " + randParagraph
		randParagraph = randomdata.Paragraph()
	}

	body := make(map[string]interface{})
	body["name"] = leagueName
	body["description"] = description
	body["publicView"] = publicView
	body["publicJoin"] = publicJoin
	body["signupStart"] = signupStart.Unix()
	body["signupEnd"] = signupEnd.Unix()
	body["leagueStart"] = leagueStart.Unix()
	body["leagueEnd"] = leagueEnd.Unix()

	println("creating league with name " + leagueName)

	return &league{
		Id:          makeApiCallAndGetId(t, body, "POST", "api/leagues", 200),
		Name:        leagueName,
		PublicView:  publicView,
		PublicJoin:  publicJoin,
		LeagueStart: &leagueStart,
		LeagueEnd:   &leagueEnd,
	}
}

func createTeam(t *testing.T, teams []*team, l *league, strength int) *team {
	teamName := randomdata.SillyName()
	fmt.Printf("Team Name: %v, Strength: %v\n", teamName, strength)
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

	threeParagraphs := randomdata.Paragraph() + "\n\n" +
		randomdata.Paragraph() + "\n\n" +
		randomdata.Paragraph() + "\n\n"

	body := make(map[string]interface{})
	body["name"] = teamName
	body["tag"] = tag
	body["description"] = threeParagraphs[0:min(499, len(threeParagraphs)-1)]

	return &team{
		Id:       makeApiCallAndGetId(t, body, "POST", "api/teams/", 200),
		LeagueId: l.Id,
		Name:     teamName,
		Tag:      tag,
		Wins:     0,
		Losses:   0,
		Strength: strength,
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
