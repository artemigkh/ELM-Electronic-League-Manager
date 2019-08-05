package lolApi

import (
	"bytes"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (api *NativeLoLTournamentApi) callTournamentServer(url string, body *bytes.Buffer) ([]byte, error) {

	//httpReq, err := http.NewRequest(
	//	"GET",
	//	fmt.Sprintf("https://americas.api.riotgames.com/lol/tournament-stub/v4/"+url),
	//	nil)
	//httpReq.Header.Set("X-Riot-Token", api.apiKey)
	//if body != nil {
	//	httpReq.Header.Set("Content-Type", "application/json")
	//}

	//
	//res, err := api.client.Do(httpReq)
	//if err != nil {
	//	fmt.Println("Failed to execute HTTP Request. Error:", err.Error())
	//	return nil, err
	//}
	//
	//if res.StatusCode != http.StatusOK {
	//	fmt.Println(fmt.Sprintf("Non 200 Status Code: %v %v", res.StatusCode, res.Status))
	//	return nil, err
	//}
	//
	//bodyBytes, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println("Failed to read body. Error:", err.Error())
	//	return nil, err
	//}

	return nil, nil
}

func (api *NativeLoLTournamentApi) RegisterTournament(leagueId int, region, tournamentName string) (providerId int, tournamentId int, err error) {
	res, err := api.r.Post(
		"https://americas.api.riotgames.com/lol/tournament-stub/v4/providers",
		req.Header{"X-Riot-Token": api.apiKey},
		req.BodyJSON(map[string]string{
			"region": region,
			"url":    "http://callbackGoesHere",
		}))
	if err != nil {
		return 0, 0, err
	} else if res.Response().StatusCode != http.StatusOK {
		return 0, 0, errors.New(res.String())
	}
	providerIdString, err := res.ToString()
	if err != nil {
		return 0, 0, err
	}

	providerId, err = strconv.Atoi(providerIdString)
	if err != nil {
		return 0, 0, err
	}

	res, err = api.r.Post(
		"https://americas.api.riotgames.com/lol/tournament-stub/v4/tournaments",
		req.Header{"X-Riot-Token": api.apiKey},
		req.BodyJSON(map[string]interface{}{
			"name":       tournamentName,
			"providerId": providerId,
		}))
	if err != nil {
		return 0, 0, err
	} else if res.Response().StatusCode != http.StatusOK {
		return 0, 0, errors.New(res.String())
	}
	tournamentIdString, err := res.ToString()
	if err != nil {
		return 0, 0, err
	}
	tournamentId, err = strconv.Atoi(tournamentIdString)
	if err != nil {
		return 0, 0, err
	}

	//todo: register one game so doesnt get cleaned up
	return providerId, tournamentId, nil
}

func (api *NativeLoLTournamentApi) CreateTournamentKey(tournamentId int, metadata string) (string, error) {
	res, err := api.r.Post(
		"https://americas.api.riotgames.com/lol/tournament-stub/v4/codes",
		req.Header{"X-Riot-Token": api.apiKey},
		req.Param{"tournamentId": tournamentId},
		req.BodyJSON(map[string]interface{}{
			"mapType":       "SUMMONERS_RIFT",
			"metadata":      metadata,
			"pickType":      "TOURNAMENT_DRAFT",
			"spectatorType": "LOBBYONLY",
			"teamSize":      5,
		}))
	if err != nil {
		return "", err
	} else if res.Response().StatusCode != http.StatusOK {
		return "", errors.New(res.String())
	}

	var tournamentStringArray []string
	if err = res.ToJSON(&tournamentStringArray); err != nil {
		return "", err
	}

	return tournamentStringArray[0], nil
}

func (api *NativeLoLTournamentApi) ForwardCompleteTournamentGame(body []byte) error {
	res, err := req.Post("http://localhost:8090/tournamentCallback", req.BodyJSON(body))
	if err != nil {
		return err
	} else if res.Response().StatusCode != 200 {
		return errors.New(res.Response().Status)
	} else {
		return nil
	}
}
