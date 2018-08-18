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
