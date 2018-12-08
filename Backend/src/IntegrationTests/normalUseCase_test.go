package IntegrationTests

import (
	"Server/scheduler"
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
		for i, manager := range l.Managers {
			newClient()
			loginAs(t, manager)
			setActiveLeague(t, l)

			manager.Team = createTeam(t, l.Teams, l, i+1)
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

		var teamIds []int
		for _, team := range l.Teams {
			teamIds = append(teamIds, int(team.Id))
		}

		s := scheduler.Scheduler{}
		est, _ := time.LoadLocation("America/New_York")
		s.InitScheduler(scheduler.DoubleRoundRobin, true, time.Hour,
			time.Date(2018, time.November, 8, 0, 0, 0, 0, est),
			time.Date(2018, time.December, 30, 0, 0, 0, 0, est),
			teamIds)
		s.AddWeeklyAvailability(time.Friday, 12+6, 0, time.Hour*2)
		s.AddWeeklyAvailability(time.Saturday, 12+4, 0, time.Hour*6)
		s.AddWeeklyAvailability(time.Sunday, 12+5, 0, time.Hour*5)
		games := s.GetSchedule()
		for _, game := range games {
			l.Games = append(l.Games, createGame(t, l, int(game.GameTime), float64(game.Team1Id), float64(game.Team2Id)))
		}

		//gameTime := moment.New()
		//gameTime.AddD
		//currTime := int(time.Now().Unix())
		//gameTime := randomdata.Number(currTime-2592000, currTime+2592000)
		//
		//for i:=0; i < 10; i++ {
		//	for j := i+1; j < 9; j++ {
		//		print(i)
		//		print(j)
		//		l.Games = append(l.Games, createGame(t, l, gameTime, l.Teams[i].Id, l.Teams[j].Id))
		//		//gameTime = randomdata.Number(currTime-2592000, currTime+2592000)
		//	}
		//}

		for _, g := range l.Games {
			checkGameAgainstRepresentation(t, g)
		}

		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	t.Run("Randomly unschedule 2 games", func(t *testing.T) {
		randomlyUnscheduleGames(t, l, 2)
		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	t.Run("Randomize result for past games and report them", func(t *testing.T) {
		currTime := float64(time.Now().Unix())
		for _, g := range l.Games {
			if g.GameTime < currTime {
				randomlyDecideAndReportGame(t, g, l.Teams)
			}
		}

		for _, g := range l.Games {
			checkGameAgainstRepresentation(t, g)
		}

		checkGamesAgainstLeagueSummary(t, l.Games)
	})

	t.Run("Check that standings are sorted correctly", checkTeamStandingsSortedProperly)
}
