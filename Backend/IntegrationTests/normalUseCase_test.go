package IntegrationTests

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	t.Run("each manager creates a team", func(t *testing.T) {
		for _, manager := range l.Managers {
			newClient()
			loginAs(t, manager)
			setActiveLeague(t, l)

			manager.Team = createTeam(t, l.Teams, l)
			manager.Team.Managers = append(manager.Team.Managers, manager)
			l.Teams = append(l.Teams, manager.Team)
			checkTeamCreated(t, manager.Team)
		}

		newClient()
		setActiveLeague(t, l)
		checkTeamsAgainstLeagueSummary(t, l.Teams)
	})

	t.Run("check that getting team manager endpoint returns correct information", func(t *testing.T) {
		newClient()
		loginAs(t, leagueOwner)
		setActiveLeague(t, l)
		checkLeagueManagersCorrect(t, l)
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

	t.Run("League Owner schedules round robin for all teams", func(t *testing.T) {
		newClient()
		loginAs(t, leagueOwner)
		setActiveLeague(t, l)

		gameTime := randomdata.Number(int(time.Now().Unix())-5184000, int(time.Now().Unix()))

		for _, m1 := range l.Teams {
			for _, m2 := range l.Teams {
				if m1.Id != m2.Id {
					l.Games = append(l.Games, createGame(t, l, gameTime, m1.Id, m2.Id))
					gameTime = randomdata.Number(int(time.Now().Unix())-5184000, int(time.Now().Unix()))
				}
			}
		}

		for _, g := range l.Games {
			checkGameAgainstRepresentation(t, g)
		}

		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	t.Run("Randomly unschedule 10 games", func(t *testing.T) {
		randomlyUnscheduleGames(t, l, 10)
		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	//t.Run("Randomize result for each game and report it", func(t *testing.T) {
	//	for _, g := range l.Games {
	//		randomlyDecideAndReportGame(t, g)
	//	}
	//
	//	for _, g := range l.Games {
	//		checkGameAgainstRepresentation(t, g)
	//	}
	//
	//	checkGamesAgainstLeagueSummary(t, l.Games)
	//})

	t.Run("Randomize result for each game and report it", func(t *testing.T) {
		for _, g := range l.Games {
			if randomdata.Boolean() {
				randomlyDecideAndReportGame(t, g)
			}
		}

		for _, g := range l.Games {
			checkGameAgainstRepresentation(t, g)
		}

		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	t.Run("Check that standings are sorted correctly", checkTeamStandingsSortedProperly)
}
