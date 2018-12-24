package IntegrationTests

import "time"

type league struct {
	Id          float64
	Name        string
	PublicView  bool
	PublicJoin  bool
	Teams       []*team
	Managers    []*user
	Players     []*player
	Games       []*game
	LeagueStart *time.Time
	LeagueEnd   *time.Time
}

type team struct {
	Id       float64
	LeagueId float64
	Name     string
	Tag      string
	Wins     float64
	Losses   float64
	Managers []*user
	Players  []*player
	Strength int
}

type user struct {
	Email    string
	Password string
	Leagues  []*league
	Team     *team
}

type player struct {
	Id             float64
	TeamId         float64
	GameIdentifier string
	Name           string
	mainRoster     bool
}

type game struct {
	Id         float64
	LeagueId   float64
	Team1Id    float64
	Team2Id    float64
	GameTime   float64
	Complete   bool
	WinnerId   float64
	ScoreTeam1 float64
	ScoreTeam2 float64
}
