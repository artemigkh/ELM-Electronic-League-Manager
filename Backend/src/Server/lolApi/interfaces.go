package lolApi

import (
	"Server/dataModel"
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
