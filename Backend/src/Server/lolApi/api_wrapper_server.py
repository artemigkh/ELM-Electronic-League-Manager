import json

import requests
from bottle import request, route, run

import cassiopeia as cass
from cassiopeia import Queue
from cassiopeia.core import Summoner, Match

cass.apply_settings("Backend/src/Server/lolApi/cass_config.json")

known_summoners = dict()

http = requests.session()


@route('/summonerId', method='GET')
def get_summoner_id():
    summoner_name = request.query.name
    summoner = Summoner(name=summoner_name)
    return {"id": summoner.id}


@route('/summonerInformation', method='GET')
def get_summoner_information():
    summoner_id = request.query.id
    if summoner_id in known_summoners:
        summoner = known_summoners[summoner_id]
    else:
        summoner = Summoner(id=summoner_id)
        known_summoners[summoner_id] = summoner

    information = {
        "gameIdentifier": summoner.name
    }

    if Queue.ranked_solo_fives in summoner.ranks and Queue.ranked_flex_fives in summoner.ranks:
        if summoner.ranks[Queue.ranked_solo_fives] > summoner.ranks[Queue.ranked_flex_fives]:
            information["rank"] = summoner.ranks[Queue.ranked_solo_fives].division.value
            information["tier"] = summoner.ranks[Queue.ranked_solo_fives].tier.value
        else:
            information["rank"] = summoner.ranks[Queue.ranked_flex_fives].division.value
            information["tier"] = summoner.ranks[Queue.ranked_flex_fives].tier.value

    elif Queue.ranked_solo_fives in summoner.ranks:
        information["rank"] = summoner.ranks[Queue.ranked_solo_fives].division.value
        information["tier"] = summoner.ranks[Queue.ranked_solo_fives].tier.value

    elif Queue.ranked_flex_fives in summoner.ranks:
        information["rank"] = summoner.ranks[Queue.ranked_flex_fives].division.value
        information["tier"] = summoner.ranks[Queue.ranked_flex_fives].tier.value

    else:
        information["rank"] = ""
        information["tier"] = ""

    return information


@route('/tournamentCallback', method='POST')
def get_summoner_id():
    body = request.json
    embedded_data = json.loads(body["metaData"])

    game_info_from_callback = {k: embedded_data[k] for k in ["gameId", "team1Id", "team2Id", "team1RefPlayerId"]}
    game_info = {**game_info_from_callback, **get_game_stats(body["gameId"])}
    if game_info_from_callback["team1RefPlayerId"] in game_info["winningTeamSummonerIds"]:
        game_info["winningTeamId"] = game_info["team1Id"]
        game_info["losingTeamId"] = game_info["team2Id"]
    else:
        game_info["winningTeamId"] = game_info["team2Id"]
        game_info["losingTeamId"] = game_info["team1Id"]

    print(game_info)
    http.post("http://localhost:8080/api/v1/lol/receiveCompletedTournamentGame", json=game_info)


@route('/gameStats', method='GET')
def get_game_stats_handler():
    return get_game_stats(request.query.id)


def get_game_stats(match_id):
    match = Match(id=match_id)

    stats = {
        "duration": match.duration.total_seconds(),
        "timestamp": match.creation.timestamp,
        "bannedChampions": [],
        "winningChampions": [],
        "losingChampions": [],
        "winningTeamSummonerIds": [],
        "winningTeamStats": {},
        "losingTeamSummonerIds": [],
        "losingTeamStats": {},
        "playerStats": []
    }

    if match.blue_team.win:
        winning_team = match.blue_team
        losing_team = match.red_team
    else:
        winning_team = match.red_team
        losing_team = match.blue_team

    stats["winningTeamStats"] = {
        "firstBlood": winning_team.first_blood,
        "firstTower": winning_team.first_tower,
        "side": winning_team.side.value
    }

    stats["losingTeamStats"] = {
        "firstBlood": losing_team.first_blood,
        "firstTower": losing_team.first_tower,
        "side": losing_team.side.value
    }

    for champion in winning_team.bans:
        if champion is not None:
            stats["bannedChampions"].append(champion.name)

    for champion in losing_team.bans:
        if champion is not None:
            stats["bannedChampions"].append(champion.name)

    for participant in match.participants:
        if participant.team.win:
            stats["winningChampions"].append(participant.champion.name)
            stats["winningTeamSummonerIds"].append(participant.summoner.id)
        else:
            stats["losingChampions"].append(participant.champion.name)
            stats["losingTeamSummonerIds"].append(participant.summoner.id)

        stats["playerStats"].append({
            "id": participant.summoner.id,
            "name": participant.summoner.name,
            "win": participant.team.win,
            "championPicked": participant.champion.name,
            "gold": participant.stats.gold_earned,
            "cs": participant.stats.total_minions_killed,
            "damage": participant.stats.total_damage_dealt_to_champions,
            "kills": participant.stats.kills,
            "deaths": participant.stats.deaths,
            "assists": participant.stats.assists,
            "wards": participant.stats.wards_placed,
        })

    return stats


run(host='localhost', port=8090, debug=True)
