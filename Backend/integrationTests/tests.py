from datetime import timedelta, datetime
import unittest

import requests

from .user import User
from .league import League
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
        print("logging in as ")
        print(user.email)
        print(user.password)
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
        league.assert_server_data_consistent(self)

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
        self.assertEqual(False, r.json()["leaguePermissions"]["createTeams"])
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
        self.assertEqual(201, r.status_code)
        return r.json()

    def ensure_team_creation_fails(self, name, tag, expected_error):
        r = self.http.post("http://localhost:8080/api/v1/teams", json={
            "name": name,
            "tag": tag
        })
        self.assertEqual(400, r.status_code)
        self.assertEqual(expected_error, r.json()["errorDescription"])

    def ensure_player_creation_fails(self, team_id, name, game_identifier, expected_error):
        r = self.http.post("http://localhost:8080/api/v1/teams/{}/players".format(team_id), json={
            "name": name,
            "gameIdentifier": game_identifier,
            "mainRoster": True
        })
        self.assertEqual(400, r.status_code)
        self.assertEqual(expected_error, r.json()["errorDescription"])

    def test_teams(self):
        # set up league
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)
        self.set_active_league(league)

        # check that can't make invalid teams
        self.ensure_team_creation_fails("", "TAG", "Name too short")
        self.ensure_team_creation_fails("name", "A", "Tag too short")
        self.ensure_team_creation_fails("name", "123456", "Tag too long")
        self.ensure_team_creation_fails("a" * 51, "TAG", "Name too long")

        # make a team and ensure can't make duplicates
        team1 = league.create_team(self, league_owner, "TEAM", "TAG")
        self.ensure_team_creation_fails("TEAM", "UNQ", "Name or tag in use")
        self.ensure_team_creation_fails("UNIQUE_NAME", "TAG", "Name or tag in use")

        # make sure can make a second team correctly
        team2 = league.create_team(self, league_owner, "TEAM2", "TAG2")

        # check that can't make invalid players
        self.ensure_player_creation_fails(team1.team_id, "a", "gameName", "Name too short")
        self.ensure_player_creation_fails(team1.team_id, "a" * 51, "gameName", "Name too long")
        self.ensure_player_creation_fails(team1.team_id, "Name", "", "Game identifier too short")
        self.ensure_player_creation_fails(team1.team_id, "Name", "a", "Game identifier too short")
        self.ensure_player_creation_fails(team1.team_id, "Name", "a" * 51, "Game identifier too long")

        # make a player and check that can't make duplicates on any team
        team1.add_player(self, True, "Name", "GameIdentifier")
        self.ensure_player_creation_fails(
            team1.team_id, "Name", "GameIdentifier", "Player game Identifier already in use")
        self.ensure_player_creation_fails(
            team2.team_id, "Name", "GameIdentifier", "Player game Identifier already in use")

        # make sure can make player with same name
        team1_player = team1.add_player(self, True, "Name", "GameIdentifier2")
        team2.add_player(self, True, "Name", "GameIdentifier3")

        self.check_all_teams(league)

        # updating and deleting
        team1_player.update_player(self, True, "Name", "GameIdentifierUpdated")
        self.check_all_teams(league)

        team1.remove_player(self, team1_player.player_id)
        self.check_all_teams(league)

    def test_league_permissions(self):
        # Create league
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)

        # create spectator and team manager, and owner of other league
        self.new_session()
        other_league_owner = User(self)
        self.login(other_league_owner)
        League(self)
        self.set_active_league(league)
        self.join_active_league()

        self.new_session()
        manager = User(self)
        self.login(manager)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(manager)
        team = league.create_team(self, manager)

        self.new_session()
        spectator = User(self)
        self.login(spectator)
        self.set_active_league(league)
        self.join_active_league()

        # check that unauthorized users get denied management permission endpoints
        for user in [spectator, manager, other_league_owner]:
            self.new_session()
            self.login(user)
            self.set_active_league(league)
            self.assertEqual(403, self.http.put("http://localhost:8080/api/v1/leagues", json={}).status_code)
            self.assertEqual(403, self.http.put("http://localhost:8080/api/v1/leagues/markdown", json={}).status_code)
            self.assertEqual(
                403, self.http.get("http://localhost:8080/api/v1/leagues/teamManagers", json={}).status_code)
            self.assertEqual(
                403, self.http.put("http://localhost:8080/api/v1/leagues/permissions/0", json={}).status_code)
            self.assertEqual(
                403, self.http.post("http://localhost:8080/api/v1/availabilities", json={}).status_code)
            self.assertEqual(
                403, self.http.delete("http://localhost:8080/api/v1/availabilities/0", json={}).status_code)
            self.assertEqual(
                403, self.http.post("http://localhost:8080/api/v1/weeklyAvailabilities", json={}).status_code)
            self.assertEqual(
                403, self.http.delete("http://localhost:8080/api/v1/weeklyAvailabilities/0", json={}).status_code)
            self.assertEqual(
                403, self.http.put("http://localhost:8080/api/v1/weeklyAvailabilities/0", json={}).status_code)
            self.assertEqual(
                403, self.http.post("http://localhost:8080/api/v1/schedule", json={}).status_code)
            if user == spectator:
                self.assertEqual(
                    403, self.http.put("http://localhost:8080/api/v1/teams/{}/permissions/0".format(team.team_id),
                                       json={}).status_code)

    def test_team_permissions(self):
        # TODO: test cant edit other leagues teams
        # Create league
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)
        self.set_active_league(league)

        # create manager that makes the team of interest
        self.new_session()
        manager_with_team = User(self)
        self.login(manager_with_team)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(manager_with_team)
        team = league.create_team(self, manager_with_team)
        player = team.add_player(self)

        # create manager that makes other team
        self.new_session()
        other_manager = User(self)
        self.login(other_manager)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(other_manager)
        league.create_team(self, other_manager)

        # create league spectator
        self.new_session()
        spectator = User(self)

        # bring league out of signup period
        self.new_session()
        self.login(league_owner)
        self.set_active_league(league)
        league.update_permissions(self, False)
        league.assert_server_data_consistent(self)

        # check that spectators and other managers get denied management permission endpoints
        for user in [spectator, other_manager]:
            self.new_session()
            self.login(user)
            self.set_active_league(league)
            self.assertEqual(403, self.http.post("http://localhost:8080/api/v1/teams", json={
                "name": "newTeamName",
                "description": "newTeamDescription",
                "tag": "NWTAG"
            }).status_code)

            self.assertEqual(403, self.http.put("http://localhost:8080/api/v1/teams/{}".format(team.team_id), json={
                "name": "newTeamName",
                "description": "newTeamDescription",
                "tag": "NWTAG"
            }).status_code)

            self.assertEqual(403,
                             self.http.delete("http://localhost:8080/api/v1/teams/{}".format(team.team_id)).status_code)

            self.assertEqual(403,
                             self.http.post("http://localhost:8080/api/v1/teams/{}/players".format(team.team_id), json={
                                 "name": "newPlayerName",
                                 "gameIdentifier": "newPlayerGameIdentifier",
                                 "mainRoster": True
                             }).status_code)

            self.assertEqual(
                403, self.http.put("http://localhost:8080/api/v1/teams/{}/players/{}".format(
                    team.team_id, player.player_id), json={
                    "name": "newPlayerName",
                    "gameIdentifier": "newPlayerGameIdentifier",
                    "mainRoster": True
                }).status_code)

            self.assertEqual(
                403, self.http.delete("http://localhost:8080/api/v1/teams/{}/players/{}".format(
                    team.team_id, player.player_id)).status_code)

    def test_games(self):
        # set up league
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)
        self.set_active_league(league)
        league.update_to_middle_of_competition_time(self)

        # create 3 teams
        team1 = league.create_team(self, league_owner, "TEAM1", "TAG1")
        team2 = league.create_team(self, league_owner, "TEAM2", "TAG2")
        team3 = league.create_team(self, league_owner, "TEAM3", "TAG3")

        # create two games correctly
        epoch_time = int(datetime.utcnow().timestamp())
        game1 = league.create_game(self, team1.team_id, team2.team_id, epoch_time)
        game2 = league.create_game(self, team2.team_id, team3.team_id, epoch_time + 3600)
        print(game1.__dict__)

        # check that can't create invalid games
        Game.create_game_expect_fail(self, team1.team_id - 1, team2.team_id, epoch_time + 1800,
                                     "A team in this game does not exist in this league")
        Game.create_game_expect_fail(self, team3.team_id, team2.team_id, epoch_time,
                                     "A team in this game already has a game starting at this time")
        Game.create_game_expect_fail(self, team3.team_id, team3.team_id, epoch_time + 1800,
                                     "The two teams in game must be different")
        Game.create_game_expect_fail(self, team1.team_id, team2.team_id, int(league.league_end.timestamp()) + 1,
                                     "This game start time is not during the league competition period")

        # check that can't reschedule with invalid parameters
        game1.reschedule_expect_fail(self, epoch_time + 3600,
                                     "A team in this game already has a game starting at this time")
        game1.reschedule_expect_fail(self, int(league.league_start.timestamp()) - 1,
                                     "This game start time is not during the league competition period")

        # check that can't report with invalid parameters
        game1.report_expect_fail(self, team3.team_id, team2.team_id, 3, 2,
                                 "The teams in this game report are not in this game")

        # can reschedule correctly
        game1.reschedule(self, epoch_time + 7200)

    def test_game_permissions(self):
        # Create league
        league_owner = User(self)
        self.login(league_owner)
        league = League(self)
        self.set_active_league(league)

        # create 3 managers that make two teams
        self.new_session()
        manager1 = User(self)
        self.login(manager1)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(manager1)
        team1 = league.create_team(self, manager1)

        self.new_session()
        manager2 = User(self)
        self.login(manager2)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(manager2)
        team2 = league.create_team(self, manager2)

        self.new_session()
        manager3 = User(self)
        self.login(manager3)
        self.set_active_league(league)
        self.join_active_league()
        league.managers.append(manager3)
        league.create_team(self, manager3)

        # create league spectator
        self.new_session()
        spectator = User(self)

        # league manager schedules one game correctly
        self.new_session()
        self.login(league_owner)
        self.set_active_league(league)
        league.update_to_middle_of_competition_time(self)
        epoch_time = int(datetime.utcnow().timestamp())
        game = league.create_game(self, team1.team_id, team2.team_id, epoch_time)

        # check that spectators and other managers get denied management permission endpoints
        for user in [spectator, manager2, manager3]:
            self.new_session()
            self.login(user)
            self.set_active_league(league)
            self.assertEqual(403, self.http.post("http://localhost:8080/api/v1/games", json={
                "team1Id": team1.team_id,
                "team2Id": team2.team_id,
                "gameTime": epoch_time
            }).status_code)

            self.assertEqual(403, self.http.delete("http://localhost:8080/api/v1/games/{}".format(game.game_id), json={
                "team1Id": team1.team_id,
                "team2Id": team2.team_id,
                "gameTime": epoch_time
            }).status_code)

            if user != manager2:
                print("user not manager2")
                self.assertEqual(403,
                                 self.http.post("http://localhost:8080/api/v1/games/{}/reschedule".format(game.game_id),
                                                json={"gameTime": epoch_time + 1800}).status_code)

                self.assertEqual(403,
                                 self.http.post("http://localhost:8080/api/v1/games/{}/report".format(game.game_id),
                                                json={
                                                    "winnerId": team1.team_id,
                                                    "loserId": team2.team_id,
                                                    "scoreTeam1": 2,
                                                    "scoreTeam2": 1
                                                }).status_code)

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
            new_team = league.create_team(self, manager)
            self.check_team(new_team)

        # check that league-wide team related information is correct
        self.new_session()
        self.login(league_owner)
        self.set_active_league(league)
        self.check_all_teams(league)
        # self.check_managers(league)

        league.update_to_middle_of_competition_time(self)

        # schedule double robin for all teams
        league.create_availability(self, league, "friday", 18, 0, 2 * 60)
        league.create_availability(self, league, "saturday", 16, 0, 6 * 60)
        league.create_availability(self, league, "sunday", 17, 0, 5 * 60)
        schedule = self.get_json_schedule("doubleroundrobin", 2, 1, 60)
        for json_game in schedule:
            game = league.create_game(
                self,
                json_game["team1"]["teamId"],
                json_game["team2"]["teamId"],
                json_game["gameTime"]
            )
            self.check_game(game, league.teams)
        self.check_all_games(league)

        # randomly decide result and report all games before current time
        for game in league.games:
            game.decide_result_and_report(self, league.teams)
            self.check_game(game, league.teams)
        self.check_all_games(league)


if __name__ == '__main__':
    unittest.main()
