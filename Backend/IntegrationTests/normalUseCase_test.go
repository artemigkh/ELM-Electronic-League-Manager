package IntegrationTests

import (
	"testing"
)

func Test_NormalUseCase(t *testing.T) {
    createRouterAndHttpClient()

    var u *user

    t.Run("create user, login, and check that logged in", func(t *testing.T) {
		u = createUser(t)
		loginAs(t, u)
		checkLoggedIn(t, u)
	})

}
