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
    def __init__(self, t, team_id, main_roster):
        self.team_id = team_id
        self.main_roster = main_roster

        player_profile = fake.simple_profile()
        self.name = player_profile["name"]
        self.game_identifier = player_profile["username"]

        r = t.http.post("http://localhost:8080/api/v1/teams/{}/players".format(team_id), json={
            "name": self.name,
            "gameIdentifier": self.game_identifier,
            "mainRoster": self.main_roster
            })

        t.assertEqual(201, r.status_code)
        self.player_id = r.json()["playerId"]


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

    def assert_equal_json(self, t, json):
        t.assertEqual(self.team_id, json["teamId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.description, json["description"])
        t.assertEqual(self.tag, json["tag"])
        t.assertEqual(self.wins, json["wins"])
        t.assertEqual(self.losses, json["losses"])

    def assert_display_equal_json(self, t, json):
        t.assertEqual(self.team_id, json["teamId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.tag, json["tag"])
