package lolApi

import (
	"Server/config"
	"Server/dataModel"
	"github.com/imroc/req"
	"net/http"
)

type SummonerInformation struct {
	GameIdentifier string `json:"gameIdentifier"`
	Rank           string `json:"rank"`
	Tier           string `json:"tier"`
}

type LoLApi interface {
	GetSummonerInformation(ids []string) map[string]*SummonerInformation
	CompletePlayerStubs(team *dataModel.LoLTeamStub) (*dataModel.LoLTeamWithRosters, error)
	GetSummonerId(name string) (string, error)
	GetMatchStats(id string) (*dataModel.LoLMatchInformation, error)
}

type lolApi struct {
	client *http.Client
}

func GetLolApiWrapper() LoLApi {
	//startWrapperServer()
	return &lolApi{
		client: &http.Client{},
	}
}

type LoLTournamentApi interface {
	RegisterTournament(leagueId int, region, tournamentName string) (providerId int, tournamentId int, err error)
	CreateTournamentKey(tournamentId int, metadata string) (string, error)
	ForwardCompleteTournamentGame(body []byte) error
}

type NativeLoLTournamentApi struct {
	r      *req.Req
	apiKey string
}

func GetLoLTournamentApi(config config.Config) LoLTournamentApi {
	//startWrapperServer()
	req.Debug = true
	return &NativeLoLTournamentApi{
		r:      req.New(),
		apiKey: config.GetLeagueOfLegendsApiKey(),
	}
}
