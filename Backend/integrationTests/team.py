import random
from datetime import timedelta, datetime

from faker import Faker
from faker.providers import internet
from faker.providers import lorem

fake = Faker()
fake.add_provider(internet)
fake.add_provider(lorem)


class Team:
    def __init__(self, t, league, strength):
        self.strength = strength
        self.managers = []

        self.name = fake.slug()
        self.description = fake.text(max_nb_chars=500)

        # Get Unique Tag
        base_tag = self.name[0:4].upper()
        suffix = 0
        self.tag = base_tag
        while self.tag in [team.tag for team in league.teams]:
            self.tag = base_tag + str(suffix)
            suffix += 1

        r = t.http.post("http://localhost:8080/api/v1/teams", json={
            "name": self.name,
            "description": self.description,
            "tag": self.tag
        })

        t.assertEqual(201, r.status_code)
        self.team_id = r.json()["teamId"]

        self.wins = 0
        self.losses = 0
        self.players = []

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
