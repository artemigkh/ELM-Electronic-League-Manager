package IntegrationTests

type league struct {
	Id         int
	Name       string
	PublicView bool
	PublicJoin bool
	Teams []team
}

type team struct {
	Id int
	LeagueId int
	Name string
	Tag string
	Wins int
	Losses int
	Users []user
	Players []player
}

type user struct {
	Id       int
	Email    string
	Password string
	Leagues  []*league
}

type player struct {
	Id int
	TeamId int
	GameIdentifier string
	Name string
	mainRoster bool
}

type game struct {
	Id int
	LeagueId int
	Team1Id int
	Team2Id int
	GameTime int
	Complete bool
	WinnerId int
	ScoreTeam1 int
	ScoreTeam2 int
}