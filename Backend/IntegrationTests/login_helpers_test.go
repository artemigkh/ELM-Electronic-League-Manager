package IntegrationTests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func loginAs(t *testing.T, u *user) {
	body := make(map[string]interface{})
	body["email"] = u.Email
	body["password"] = u.Password

	makeApiCall(t, body, "POST", "login", 200)
}

func checkLoggedIn(t *testing.T, u *user) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET", "api/users/profile", 200)
	assert.Equal(t, responseMap["email"].(string), u.Email)
}

func logout(t *testing.T) {
	makeApiCall(t, nil, "POST", "logout", 200)
}

func checkLoggedOut(t *testing.T) {
	responseMap := makeApiCallAndGetMap(t, nil, "GET", "api/users/profile", 403)
	assert.Equal(t, responseMap["error"].(string), "notLoggedIn")
}
