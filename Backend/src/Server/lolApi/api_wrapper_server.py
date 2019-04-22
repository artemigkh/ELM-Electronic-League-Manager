from bottle import request, route, run

import cassiopeia as cass
from cassiopeia import Queue
from cassiopeia.core import Summoner

cass.apply_settings("Backend/src/Server/lolApi/cass_config.json")


@route('/summonerId', method='GET')
def get_summoner_id():
    summoner_name = request.query.name
    summoner = Summoner(name=summoner_name)
    return {"id:": summoner.id}


@route('/summonerInformation', method='GET')
def get_summoner_information():
    summoner_id = request.query.id
    summoner = Summoner(id=summoner_id)
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


run(host='localhost', port=8090, debug=True)
