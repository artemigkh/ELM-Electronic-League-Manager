from datetime import datetime, timedelta
from pytz import timezone


class Availability:
    def __init__(self, t, league, weekday, hour, minute, duration_minutes):
        self.start_time = int(league.league_start.timestamp())
        self.end_time = int(league.league_end.timestamp())

        self.weekday = weekday
        self.timezone = timezone('US/Eastern').utcoffset(datetime.utcnow(), is_dst=False).seconds
        self.hour = hour
        self.minute = minute
        self.duration_minutes = duration_minutes

        r = t.http.post("http://localhost:8080/api/v1/weeklyAvailabilities", json={
            "startTime": self.start_time,
            "endTime": self.end_time,
            "weekday": self.weekday,
            "timezone": self.timezone,
            "hour": self.hour,
            "minute": self.minute,
            "duration": self.duration_minutes
        })

        if r.status_code != 201:
            print(r.json())
        t.assertEqual(201, r.status_code)
        self.availability_id = r.json()["availabilityId"]

    def assert_equal_json(self, t, json):
        t.assertEqual(self.league_id, json["leagueId"])
        t.assertEqual(self.name, json["name"])
        t.assertEqual(self.description, json["description"])
        t.assertEqual(self.game, json["game"])
        t.assertEqual(int(self.signup_start.timestamp()), json["signupStart"])
        t.assertEqual(int(self.signup_end.timestamp()), json["signupEnd"])
        t.assertEqual(int(self.league_start.timestamp()), json["leagueStart"])
        t.assertEqual(int(self.league_end.timestamp()), json["leagueEnd"])
