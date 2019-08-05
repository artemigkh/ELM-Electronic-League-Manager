import random
from datetime import timedelta, datetime

from faker import Faker
from faker.providers import internet
from faker.providers import profile
from faker.providers import lorem

fake = Faker()
fake.add_provider(internet)
fake.add_provider(lorem)
fake.add_provider(profile)


class Player:
    def __init__(self, t, team_id, main_roster, name=None, game_identifier=None):
        self.team_id = team_id
        self.main_roster = main_roster

        player_profile = fake.simple_profile()
        if name is None:
            self.name = player_profile["name"]
        else:
            self.name = name

        if game_identifier is None:
            self.game_identifier = player_profile["username"]
        else:
            self.game_identifier = game_identifier
        r = t.http.post("http://localhost:8080/api/v1/teams/{}/players".format(team_id), json={
            "name": self.name,
            "gameIdentifier": self.game_identifier,
            "mainRoster": self.main_roster
        })

        t.assertEqual(201, r.status_code)
        self.player_id = r.json()["playerId"]

    def update_player(self, t, main_roster, name, game_identifier):
        self.main_roster = main_roster
        self.name = name
        self.game_identifier = game_identifier
        r = t.http.put("http://localhost:8080/api/v1/teams/{}/players/{}".
                       format(self.team_id, self.player_id),
                       json={
                           "name": self.name,
                           "gameIdentifier": self.game_identifier,
                           "mainRoster": self.main_roster
                       })

        t.assertEqual(200, r.status_code)

    def assert_equal_json(self, t, json):
        t.assertEqual(self.player_id, json["playerId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.game_identifier, json["gameIdentifier"])
        t.assertEqual(self.main_roster, json["mainRoster"])


class Team:
    def __init__(self, t, league, manager, strength, name=None, tag=None):
        self.strength = strength
        self.managers = [manager]
        self.players = []

        if name is None:
            self.name = fake.slug()
        else:
            self.name = name

        self.description = fake.text(max_nb_chars=500)

        if tag is None:
            # Get Unique Tag
            base_tag = self.name[0:4].upper()
            suffix = 0
            self.tag = base_tag
            while self.tag in [team.tag for team in league.teams]:
                self.tag = base_tag + str(suffix)
                suffix += 1
        else:
            self.tag = tag
        r = t.http.post("http://localhost:8080/api/v1/teams", json={
            "name": self.name,
            "description": self.description,
            "tag": self.tag
        })

        t.assertEqual(201, r.status_code)
        self.team_id = r.json()["teamId"]

        self.wins = 0
        self.losses = 0

        # Add 5 main roster players and 2 substitutes
        self.players = []
        for i in range(7):
            self.players.append(Player(t, self.team_id, i < 5))

    def get_player(self, player_id):
        return next((p for p in self.players if p.player_id == player_id), None)

    def add_player(self, t, main_roster=True, name=None, game_identifier=None):
        new_player = Player(t, self.team_id, main_roster, name, game_identifier)
        self.players.append(new_player)
        return new_player

    def remove_player(self, t, player_id):
        to_delete = self.get_player(player_id)
        r = t.http.delete("http://localhost:8080/api/v1/teams/{}/players/{}".format(self.team_id, player_id))
        t.assertEqual(200, r.status_code)
        self.players.remove(to_delete)

    def assert_equal_json(self, t, json):
        t.assertEqual(self.team_id, json["teamId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.description, json["description"])
        t.assertEqual(self.tag, json["tag"])
        t.assertEqual(self.wins, json["wins"])
        t.assertEqual(self.losses, json["losses"])

        def get_json_player(player_id):
            return next((p for p in json["players"] if p["playerId"] == player_id), None)
        for player in self.players:
            t.assertIsNotNone(get_json_player(player.player_id))

        for json_player in json["players"]:
            player = self.get_player(json_player["playerId"])
            player.assert_equal_json(t, json_player)

    def assert_display_equal_json(self, t, json):
        t.assertEqual(self.team_id, json["teamId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.tag, json["tag"])
