import random
import unittest

import requests

from .user import User
from .league import League
from .team import Team
from .availability import Availability
from .game import Game


class TestElmApi(unittest.TestCase):
    def __init__(self, *args, **kwargs):
        super(TestElmApi, self).__init__(*args, **kwargs)
        self.http = requests.session()

    def new_session(self):
        self.http.close()
        del self.http
        self.http = requests.session()

    def setUp(self):
        self.new_session()

    def tearDown(self):
        self.http.close()

    def login(self, user):
        # Make request to login endpoint
        r = self.http.post("http://localhost:8080/login", json={
            "email": user.email,
            "password": user.password
        })
        self.assertEqual(200, r.status_code)

        # Get profile to see if matches with who just logged in
        r = self.http.get("http://localhost:8080/api/v1/users")

        self.assertEqual(200, r.status_code)
        self.assertEqual(user.email, r.json()["email"])

    def logout(self):
        # Make request to logout endpoint
        r = self.http.post("http://localhost:8080/logout")
        self.assertEqual(200, r.status_code)

        # Get profile to make sure logged out
        r = self.http.get("http://localhost:8080/api/v1/users")

        self.assertEqual(403, r.status_code)

    def set_active_league(self, league):
        # Make request to set active league endpoint
        r = self.http.post("http://localhost:8080/api/v1/leagues/setActiveLeague/{}".
                           format(league.league_id))
        self.assertEqual(200, r.status_code)

        # Get league info to make sure correct league set
        r = self.http.get("http://localhost:8080/api/v1/leagues")
        self.assertEqual(200, r.status_code)
        league.assert_equal_json(self, r.json())

    def join_active_league(self):
        # Check that no permissions before join
        r = self.http.get("http://localhost:8080/api/v1/users/leaguePermissions")
        self.assertEqual(200, r.status_code)
        self.assertEqual(False, r.json()["leaguePermissions"]["administrator"])
        self.assertEqual(False, r.json()["leaguePermissions"]["createTeams"])
        self.assertEqual(False, r.json()["leaguePermissions"]["editTeams"])
        self.assertEqual(False, r.json()["leaguePermissions"]["editGames"])

        # Make request to join active league endpoint
        r = self.http.post("http://localhost:8080/api/v1/leagues/join")
        self.assertEqual(200, r.status_code)

        # Check that only create teams permission gets set after join
        r = self.http.get("http://localhost:8080/api/v1/users/leaguePermissions")
        self.assertEqual(200, r.status_code)
        self.assertEqual(False, r.json()["leaguePermissions"]["administrator"])
        self.assertEqual(True, r.json()["leaguePermissions"]["createTeams"])
        self.assertEqual(False, r.json()["leaguePermissions"]["editTeams"])
        self.assertEqual(False, r.json()["leaguePermissions"]["editGames"])

    def check_team(self, team):
        r = self.http.get("http://localhost:8080/api/v1/teams/{}".format(team.team_id))
        self.assertEqual(200, r.status_code)
        team.assert_equal_json(self, r.json())

    def check_game(self, game, teams):
        r = self.http.get("http://localhost:8080/api/v1/games/{}".format(game.game_id))
        self.assertEqual(200, r.status_code)
        game.assert_equal_json(self, r.json(), teams)

    def check_all_games(self, league):
        r = self.http.get("http://localhost:8080/api/v1/games")
        self.assertEqual(200, r.status_code)
        league.assert_games_equal_json(self, r.json())

    def check_all_teams(self, league):
        r = self.http.get("http://localhost:8080/api/v1/teams")
        self.assertEqual(200, r.status_code)
        league.assert_teams_equal_json(self, r.json())

    def check_managers(self, league):
        r = self.http.get("http://localhost:8080/api/v1/leagues/teamManagers")
        self.assertEqual(200, r.status_code)
        league.assert_managers_equal_json(self, r.json())

    def get_json_schedule(self, tournament_type, rounds_per_week, concurrent_game_num, game_duration_minutes):
        r = self.http.post("http://localhost:8080/api/v1/schedule", json={
                "tournamentType": tournament_type,
                "roundsPerWeek": rounds_per_week,
                "concurrentGameNum": concurrent_game_num,
                "gameDuration": game_duration_minutes
        })
        self.assertEqual(200, r.status_code)
        return r.json()

    def test_normalUseCase(self):
        # create league owner
        league_owner = User(self)
        print("Created league owner with \nEmail: {}\nPassword: {}".format(
            league_owner.email, league_owner.password))
        self.login(league_owner)

        # check logout works
        self.logout()
        self.login(league_owner)

        # check create league and state works
        league = League(self)
        print("Created league with id " + str(league.league_id))
        self.set_active_league(league)

        # create 10 independent managers that join league
        for _ in range(10):
            self.new_session()
            manager = User(self)
            self.login(manager)
            self.set_active_league(league)
            self.join_active_league()
            league.managers.append(manager)
        self.assertEqual(10, len(league.managers))

        # each manager independently creates a team
        for manager in league.managers:
            self.new_session()
            self.login(manager)
            self.set_active_league(league)

            new_team = Team(self, league, random.randint(0, 100))
            new_team.managers.append(manager)
            league.teams.append(new_team)
            self.check_team(new_team)

        # check that league-wide team related information is correct
        self.new_session()
        self.login(league_owner)
        self.set_active_league(league)
        self.check_all_teams(league)
        # self.check_managers(league)

        # schedule double robin for all teams
        league.availabilities.append(Availability(self, league, "friday", 18, 0, 2 * 60))
        league.availabilities.append(Availability(self, league, "saturday", 16, 0, 6 * 60))
        league.availabilities.append(Availability(self, league, "sunday", 17, 0, 5 * 60))
        schedule = self.get_json_schedule("doubleroundrobin", 2, 1, 60)
        for json_game in schedule:
            game = Game(
                self,
                json_game["team1"]["teamId"],
                json_game["team2"]["teamId"],
                json_game["gameTime"]
            )
            league.games.append(game)
            self.check_game(game, league.teams)
        self.check_all_games(league)

        # randomly decide result and report all games before current time
        for game in league.games:
            game.decide_result_and_report(self, league.teams)
            self.check_game(game, league.teams)
        self.check_all_games(league)

    def test_Games(self):
        # set up league and 2 teams
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)
        self.set_active_league(league)
        team1 = Team(self, league, random.randint(0, 100))
        team2 = Team(self, league, random.randint(0, 100))


if __name__ == '__main__':
    unittest.main()
