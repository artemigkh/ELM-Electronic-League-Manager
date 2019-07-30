package dataModel

type LeagueOfLegendsDAO interface {
	CreateLoLPlayer(leagueId, teamId int, externalId string, playerInfo LoLPlayerCore) (int, error)
	CreateLoLTeamWithPlayers(leagueId, userId int, teamInfo TeamCore, players []*LoLPlayerCore, iconSmall, iconLarge string) (int, error)
	UpdateLoLPlayer(playerId int, externalId string, playerInfo LoLPlayerCore) error

	GetLoLTeamStub(teamId int) (*LoLTeamStub, error)
	GetAllLoLTeamStubInLeague(leagueId int) ([]*LoLTeamStub, error)

	ReportEndGameStats(leagueId, gameId int, match *LoLMatchInformation) error
	GetPlayerStats(leagueId int) ([]*LoLPlayerStats, error)
	GetTeamStats(leagueId int) ([]*LoLTeamStats, error)
	GetChampionStats(leagueId int) ([]*LoLChampionStats, error)
}

type LoLTeamWithPlayersCore struct {
	Team    TeamCore         `json:"team"`
	Icon    string           `json:"icon"`
	Players []*LoLPlayerCore `json:"players"`
}

func (team *LoLTeamWithPlayersCore) Validate(leagueId int, teamDao TeamDAO) (bool, string, error) {
	//TODO: validate all players within team as well
	return true, "", nil
}

type LoLPlayer struct {
	PlayerId       int    `json:"playerId"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
	Position       string `json:"position"`
	Rank           string `json:"rank"`
	Tier           string `json:"tier"`
}

type LoLPlayerStub struct {
	PlayerId   int
	ExternalId string
	MainRoster bool
	Position   string
}

type LoLTeamStub struct {
	TeamId           int
	Name             string
	Description      string
	Tag              string
	IconSmall        string
	IconLarge        string
	Wins             int
	Losses           int
	MainRoster       []*LoLPlayerStub
	SubstituteRoster []*LoLPlayerStub
}

type LoLPlayerStats struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	TeamId          int     `json:"teamId"`
	AverageDuration float64 `json:"averageDuration"`
	GoldPerMinute   float64 `json:"goldPerMinute"`
	CsPerMinute     float64 `json:"csPerMinute"`
	DamagePerMinute float64 `json:"damagePerMinute"`
	AverageKills    float64 `json:"averageKills"`
	AverageDeaths   float64 `json:"averageDeaths"`
	AverageAssists  float64 `json:"averageAssists"`
	AverageKda      float64 `json:"averageKda"`
	AverageWards    float64 `json:"averageWards"`
}

type LoLTeamStats struct {
	Id                 string  `json:"id"`
	AverageDuration    float64 `json:"averageDuration"`
	NumberFirstBloods  int     `json:"numberFirstBloods"`
	NumberFirstTurrets int     `json:"numberFirstTurrets"`
	AverageKda         float64 `json:"averageKda"`
	AverageWards       float64 `json:"averageWards"`
	AverageActionScore float64 `json:"averageActionScore"`
	GoldPerMinute      float64 `json:"goldPerMinute"`
	CsPerMinute        float64 `json:"csPerMinute"`
}

type LoLChampionStats struct {
	Name    string  `json:"name"`
	Bans    int     `json:"bans"`
	Picks   int     `json:"picks"`
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	Winrate float64 `json:"winrate"`
}

type LoLMatchTeamStats struct {
	FirstBlood bool `json:"firstBlood"`
	FirstTower bool `json:"firstTower"`
	Side       int  `json:"side"`
}

type LoLMatchPlayerStats struct {
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

type LoLMatchInformation struct {
	GameId                 string                `json:"gameId"`
	Duration               float64               `json:"duration"`
	Timestamp              int                   `json:"timestamp"`
	Team1Id                int                   `json:"team1Id"`
	Team2Id                int                   `json:"team2Id"`
	WinningTeamId          int                   `json:"winningTeamId"`
	LosingTeamId           int                   `json:"losingTeamId"`
	BannedChampions        []string              `json:"bannedChampions"`
	WinningChampions       []string              `json:"winningChampions"`
	LosingChampions        []string              `json:"losingChampions"`
	WinningTeamSummonerIds []string              `json:"winningTeamSummonerIds"`
	LosingTeamSummonerIds  []string              `json:"losingTeamSummonerIds"`
	WinningTeamStats       LoLMatchTeamStats     `json:"winningTeamStats"`
	LosingTeamStats        LoLMatchTeamStats     `json:"losingTeamStats"`
	PlayerStats            []LoLMatchPlayerStats `json:"playerStats"`
}

type LoLPlayerCore struct {
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
	Position       string `json:"position"`
	ExternalId     string `json:"externalId"`
}

func (player *LoLPlayerCore) validate(leagueId, teamId, playerId int, leagueOfLegendsDAO LeagueOfLegendsDAO) (bool, string, error) {
	return validate(
		validateGameIdentifier(player.GameIdentifier),
		player.uniqueness(teamId, playerId, leagueOfLegendsDAO))
}

func (player *LoLPlayerCore) ValidateNew(leagueId, teamId int, leagueOfLegendsDAO LeagueOfLegendsDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, 0, leagueOfLegendsDAO)
}

func (player *LoLPlayerCore) ValidateEdit(leagueId, teamId, playerId int, leagueOfLegendsDAO LeagueOfLegendsDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, playerId, leagueOfLegendsDAO)
}

func (player *LoLPlayerCore) uniqueness(leagueId, playerId int, leagueOfLegendsDAO LeagueOfLegendsDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		teams, err := leagueOfLegendsDAO.GetAllLoLTeamStubInLeague(leagueId)
		if err != nil {
			errorDest = &err
			return false
		}

		for _, team := range teams {
			for _, existingPlayer := range append(team.MainRoster, team.SubstituteRoster...) {
				if existingPlayer.ExternalId == player.ExternalId && existingPlayer.PlayerId != playerId {
					*problemDest = ExternalIdentifierInUse
					return false
				}
			}
		}
		return true
	}
}
