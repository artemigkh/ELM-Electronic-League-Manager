package lolApi

import (
	"Server/databaseAccess"
	"net/http"
)

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
	GameId                 string        `json:"gameId"`
	Duration               float64       `json:"duration"`
	Timestamp              int           `json:"timestamp"`
	Team1Id                int           `json:"team1Id"`
	Team2Id                int           `json:"team2Id"`
	WinningTeamId          int           `json:"winningTeamId"`
	LosingTeamId           int           `json:"losingTeamId"`
	BannedChampions        []string      `json:"bannedChampions"`
	WinningChampions       []string      `json:"winningChampions"`
	LosingChampions        []string      `json:"losingChampions"`
	WinningTeamSummonerIds []string      `json:"winningTeamSummonerIds"`
	LosingTeamSummonerIds  []string      `json:"losingTeamSummonerIds"`
	WinningTeamStats       TeamStats     `json:"winningTeamStats"`
	LosingTeamStats        TeamStats     `json:"losingTeamStats"`
	PlayerStats            []PlayerStats `json:"playerStats"`
}

type LoLApi interface {
	GetSummonerInformation(ids []string) map[string]*SummonerInformation
	CompletePlayerStubs(team *databaseAccess.LoLTeamStub) (*databaseAccess.LoLTeamWithRosters, error)
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
