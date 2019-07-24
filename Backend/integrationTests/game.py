import random
from datetime import datetime


class Game:
    def __init__(self, t, team1_id, team2_id, game_time):
        self.team1_id = team1_id
        self.team2_id = team2_id
        self.game_time = game_time

        self.complete = False
        self.winner_id = -1
        self.loser_id = -1
        self.score_team1 = 0
        self.score_team2 = 0

        r = t.http.post("http://localhost:8080/api/v1/games", json={
            "team1Id": self.team1_id,
            "team2Id": self.team2_id,
            "gameTime": self.game_time
        })

        t.assertEqual(201, r.status_code)
        self.game_id = r.json()["gameId"]

    def decide_result_and_report(self, t, teams):
        if self.game_time < int(datetime.utcnow().timestamp()):
            team1 = next((t for t in teams if t.team_id == self.team1_id), None)
            team2 = next((t for t in teams if t.team_id == self.team2_id), None)

            while self.score_team1 < 3 and self.score_team2 < 3:
                if random.randint(0, team1.strength) > random.randint(0, team2.strength):
                    self.score_team1 += 1
                else:
                    self.score_team2 += 1

            if self.score_team1 > self.score_team2:
                self.winner_id = team1.team_id
                self.loser_id = team2.team_id
            else:
                self.winner_id = team2.team_id
                self.loser_id = team1.team_id

            self.complete = True

            r = t.http.post("http://localhost:8080/api/v1/games/{}/report".format(self.game_id), json={
                "winnerId": self.winner_id,
                "loserId": self.loser_id,
                "scoreTeam1": self.score_team1,
                "scoreTeam2": self.score_team2
            })

            t.assertEqual(200, r.status_code)

    def assert_equal_json(self, t, json, teams):
        t.assertEqual(self.game_time, json["gameTime"])
        t.assertEqual(self.complete, json["complete"])
        t.assertEqual(self.winner_id, json["winnerId"])
        t.assertEqual(self.loser_id, json["loserId"])
        t.assertEqual(self.score_team1, json["scoreTeam1"])
        t.assertEqual(self.score_team2, json["scoreTeam2"])

        team1 = next((t for t in teams if t.team_id == json["team1"]["teamId"]), None)
        team1.assert_display_equal_json(t, json["team1"])
        team2 = next((t for t in teams if t.team_id == json["team2"]["teamId"]), None)
        team2.assert_display_equal_json(t, json["team2"])
