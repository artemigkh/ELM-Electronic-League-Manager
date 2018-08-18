package IntegrationTests

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_NormalUseCase(t *testing.T) {
	createRouterAndHttpClient()

	var leagueOwner *user
	var l *league

	t.Run("create user, login, and check that logged in", func(t *testing.T) {
		leagueOwner = createUser(t)
		loginAs(t, leagueOwner)
		checkLoggedIn(t, leagueOwner)
	})

	t.Run("logout and check that can't get profile", func(t *testing.T) {
		logout(t)
		checkLoggedOut(t)
	})

	t.Run("ensure can't create leagued logged out, login, then create league", func(t *testing.T) {
		checkCantMakeLeagueLoggedOut(t)
		loginAs(t, leagueOwner)
		l = createLeague(t, true, true)
		checkCantGetLeagueNoActiveLeague(t)
		setActiveLeague(t, l)
		checkLeagueSelected(t, l)
	})

	t.Run("reset client and ensure no league is active", func(t *testing.T) {
		newClient()
		checkCantGetLeagueNoActiveLeague(t)
	})

	t.Run("10 Managers created and join league", func(t *testing.T) {
		newClient()
		for i := 0; i < 10; i++ {
			u := createUser(t)
			loginAs(t, u)
			setActiveLeague(t, l)
			joinLeague(t)
			l.Managers = append(l.Managers, u)
		}
		assert.Equal(t, len(l.Managers), 10)
	})

	//TODO: check that all 10 are registered as Managers in the league

	t.Run("each manager creates a team", func(t *testing.T) {
		for _, manager := range l.Managers {
			newClient()
			loginAs(t, manager)
			setActiveLeague(t, l)

			manager.Team = createTeam(t, l.Teams, l)
			l.Teams = append(l.Teams, manager.Team)
			checkTeamCreated(t, manager.Team)
		}

		newClient()
		setActiveLeague(t, l)
		checkTeamsAgainstLeagueSummary(t, l.Teams)
	})

	t.Run("each manager adds 5 main roster players and 2 subs", func(t *testing.T) {
		for _, manager := range l.Managers {
			newClient()
			loginAs(t, manager)
			setActiveLeague(t, l)

			for i := 0; i < 5; i++ {
				addPlayerToTeam(t, manager.Team, l, true)
			}
			for i := 0; i < 2; i++ {
				addPlayerToTeam(t, manager.Team, l, false)
			}
		}

		for _, m := range l.Teams {
			newClient()
			setActiveLeague(t, l)
			checkTeamAgainstRepresentation(t, m)
		}
	})
}
