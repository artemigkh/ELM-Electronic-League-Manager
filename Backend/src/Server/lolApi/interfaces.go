package lolApi

import "net/http"

type SummonerInformation struct {
	GameIdentifier string `json:"gameIdentifier"`
	Rank           string `json:"rank"`
	Tier           string `json:"tier"`
}

type LoLApi interface {
	GetSummonerInformation(ids []string) map[string]*SummonerInformation
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
