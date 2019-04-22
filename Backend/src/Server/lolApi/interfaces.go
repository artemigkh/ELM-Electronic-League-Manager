package lolApi

import "net/http"

type SummonerInformation struct {
	GameIdentifier string `json:"gameIdentifier"`
	Rank           string `json:"rank"`
	Tier           string `json:"tier"`
}

type TeamStats struct {
	FirstBlood bool `json:"firstBlood"`
	FirstTower bool `json:"firstTower"`
	Side       int  `json:"side"`
}

type PlayerStats struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	ChampionPicked string  `json:"championPicked"`
	Gold           float64 `json:"gold"`
	Cs             float64 `json:"cs"`
	Damage         float64 `json:"damage"`
	Kills          float64 `json:"kills"`
	Deaths         float64 `json:"deaths"`
	Assists        float64 `json:"assists"`
	Wards          float64 `json:"wards"`
	Win            bool    `json:"win"`
}

type MatchInformation struct {
	Duration         float64       `json:"duration"`
	Timestamp        int           `json:"timestamp"`
	BannedChampions  []string      `json:"bannedChampions"`
	WinningChampions []string      `json:"winningChampions"`
	LosingChampions  []string      `json:"losingChampions"`
	WinningTeamIds   []string      `json:"winningTeamIds"`
	LosingTeamIds    []string      `json:"losingTeamIds"`
	WinningTeamStats TeamStats     `json:"winningTeamStats"`
	LosingTeamStats  TeamStats     `json:"losingTeamStats"`
	PlayerStats      []PlayerStats `json:"playerStats"`
}

type LoLApi interface {
	GetSummonerInformation(ids []string) map[string]*SummonerInformation
	GetSummonerId(name string) (string, error)
	GetMatchStats(id string) (*MatchInformation, error)
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
