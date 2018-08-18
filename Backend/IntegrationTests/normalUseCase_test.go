package IntegrationTests

import (
	"testing"
)

func Test_NormalUseCase(t *testing.T) {
	createRouterAndHttpClient()

	var u *user
	var l *league

	t.Run("create user, login, and check that logged in", func(t *testing.T) {
		u = createUser(t)
		loginAs(t, u)
		checkLoggedIn(t, u)
	})

	t.Run("logout and check that can't get profile", func(t *testing.T) {
		logout(t)
		checkLoggedOut(t)
	})

	t.Run("ensure can't create leagued logged out, login, then create league", func(t *testing.T) {
		checkCantMakeLeagueLoggedOut(t)
		loginAs(t, u)
		l = createLeague(t, true, true)
		checkCantGetLeagueNoActiveLeague(t)
		setActiveLeague(t, l)
		checkLeagueSelected(t, l)
	})
}
