package IntegrationTests

import "testing"

func loginAs(t *testing.T, u *user) {
	body := make(map[string]interface{})
	body["email"] = u.Email
	body["password"] = u.Password

	makeApiCall(t, body, "POST", "login", 200)
}

func checkLoggedIn(t *testing.T, u *user) {
	println(makeApiCallAndGetId(t, nil, "GET", "api/users/profile", 200))
}
