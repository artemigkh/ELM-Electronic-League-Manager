package databaseAccess

import (
	"github.com/artemigkh/GoLang-LeagueOfLegendsAPIV4Framework"
)

type ChampionInfo struct {
	ChampionName string  `json:"championName"`
	NumGames     int     `json:"numGames"`
	NumBans      int     `json:"numBans"`
	NumWins      int     `json:"numWins"`
	NumLosses    int     `json:"numLosses"`
	Winrate      float64 `json:"winrate"`
}

type TeamInfo struct {
	TeamId    int     `json:"teamId"`
	TeamName  string  `json:"teamName"`
	TeamTag   string  `json:"teamTag"`
	StatName  string  `json:"statName"`
	StatValue float64 `json:"statValue"`
	NumGames  int     `json:"numGames"`
}

type PlayerInfo struct {
	PlayerName string  `json:"playerName"`
	StatName   string  `json:"statName"`
	StatValue  float64 `json:"statValue"`
	NumGames   int     `json:"numGames"`
}

type PlayerWardInfo struct {
	PlayerName   string `json:"playerName"`
	NumGames     int    `json:"numGames"`
	VisionWards  int    `json:"visionWards"`
	ControlWards int    `json:"controlWards"`
}

type TopPerformers struct {
	ChampionInfo []ChampionInfo `json:"championInfo"`

	// Teams
	HighestTeamKDA   []TeamInfo `json:"highestTeamKda"`
	ShortestGames    []TeamInfo `json:"shortestGames"`
	MostFirstBloods  []TeamInfo `json:"mostFirstBloods"`
	MostFirstTurrets []TeamInfo `json:"mostFirstTurrets"`

	// Players
	HighestDPM       []PlayerInfo     `json:"highestDpm"`
	HighestGPM       []PlayerInfo     `json:"higheestGpm"`
	HighestCSPM      []PlayerInfo     `json:"highestCspm"`
	HighestPlayerKDA []PlayerInfo     `json:"highestPlayerKda"`
	MostKills        []PlayerInfo     `json:"mostKills"`
	MostDeaths       []PlayerInfo     `json:"mostDeaths"`
	MostAssists      []PlayerInfo     `json:"mostAssists"`
	MostWardsPlaced  []PlayerWardInfo `json:"mostVisionWards"`
}

func getChampionInformation(leagueId int) ([]ChampionInfo, error) {
	rows, err := psql.Select("name", "picks", "wins",
		"picks - wins as losses", "bans", "wins / picks as winrate").
		From("championStats").
		Where("leagueId = ?", leagueId).
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var champions []ChampionInfo
	var champion ChampionInfo
	for rows.Next() {
		err := rows.Scan(&champion.ChampionName, &champion.NumGames,
			&champion.NumWins, &champion.NumLosses, &champion.NumBans,
			&champion.Winrate)
		if err != nil {
			return nil, err
		}
		champions = append(champions, champion)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return champions, nil
}

func getTeamInfo(leagueId int, column, label string) ([]TeamInfo, error) {
	rows, err := psql.Select("teamId", "numGames", "name", "tag", column).
		From("teamStats").Join("teams ON teamStats.teamId = teams.id").
		Where("teamStats.leagueId = ?", leagueId).
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamInfo []TeamInfo
	var teamInfoEntry TeamInfo
	for rows.Next() {
		err := rows.Scan(teamInfoEntry.TeamId, teamInfoEntry.NumGames, teamInfoEntry.TeamName,
			teamInfoEntry.TeamTag, teamInfoEntry.StatValue)
		if err != nil {
			return nil, err
		}
		teamInfoEntry.StatName = label
		teamInfo = append(teamInfo, teamInfoEntry)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return teamInfo, nil
}

func getHighestTeamKDA(leagueId int) ([]TeamInfo, error) {
	return nil, nil
}

func getShortestGames(leagueId int) ([]TeamInfo, error) {
	return nil, nil
}

func getMostFirstBloods(leagueId int) ([]TeamInfo, error) {
	return nil, nil
}

func getMostFirstTurrets(leagueId int) ([]TeamInfo, error) {
	return nil, nil
}

func getHighestDPM(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}

func getHighestGPM(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}
func getHighestCSPM(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}
func getHighestPlayerKDA(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}
func getMostKills(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}

func getMostDeaths(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}

func getMostAssists(leagueId int) ([]PlayerInfo, error) {
	return nil, nil
}

func getMostWardsPlaced(leagueId int) ([]PlayerWardInfo, error) {
	return nil, nil
}

type PgLeagueOfLegendsDAO struct{}

func (d *PgLeagueOfLegendsDAO) GetTopPerformers(leagueId int) (*TopPerformers, error) {
	championInfo, err := getChampionInformation(leagueId)
	if err != nil {
		return nil, err
	}

	highestTeamKDA, err := getHighestTeamKDA(leagueId)
	if err != nil {
		return nil, err
	}

	shortestGames, err := getShortestGames(leagueId)
	if err != nil {
		return nil, err
	}

	mostFirstBloods, err := getMostFirstBloods(leagueId)
	if err != nil {
		return nil, err
	}

	mostFirstTurrets, err := getMostFirstTurrets(leagueId)
	if err != nil {
		return nil, err
	}

	highestDPM, err := getHighestDPM(leagueId)
	if err != nil {
		return nil, err
	}

	highestGPM, err := getHighestGPM(leagueId)
	if err != nil {
		return nil, err
	}

	highestCSPM, err := getHighestCSPM(leagueId)
	if err != nil {
		return nil, err
	}

	highestPlayerKDA, err := getHighestPlayerKDA(leagueId)
	if err != nil {
		return nil, err
	}

	mostKills, err := getMostKills(leagueId)
	if err != nil {
		return nil, err
	}

	mostDeaths, err := getMostDeaths(leagueId)
	if err != nil {
		return nil, err
	}

	mostAssists, err := getMostAssists(leagueId)
	if err != nil {
		return nil, err
	}

	mostWardsPlaced, err := getMostWardsPlaced(leagueId)
	if err != nil {
		return nil, err
	}

	return &TopPerformers{
		ChampionInfo:     championInfo,
		HighestTeamKDA:   highestTeamKDA,
		ShortestGames:    shortestGames,
		MostFirstBloods:  mostFirstBloods,
		MostFirstTurrets: mostFirstTurrets,
		HighestDPM:       highestDPM,
		HighestGPM:       highestGPM,
		HighestCSPM:      highestCSPM,
		HighestPlayerKDA: highestPlayerKDA,
		MostKills:        mostKills,
		MostDeaths:       mostDeaths,
		MostAssists:      mostAssists,
		MostWardsPlaced:  mostWardsPlaced,
	}, nil
}

func (d *PgLeagueOfLegendsDAO) reportEndGameStats(match GoLang_LeagueOfLegendsAPIV4Framework.Match) error {
	for _, summ := range match.Summoners() {
		var (
			id              int
			numGames        int
			goldPerMinute   float64
			csPerMinute     float64
			damagePerMinute float64
			kills           float64
			deaths          float64
			assists         float64
			visionWards     float64
			controlWards    float64
		)
		err := psql.Select("playerId", "numGames", "goldPerMinute", "csPerMinute",
			"damagePerMinute", "kills", "deaths", "assists", "visionWards", "controlWards").
			From("playerStats").Join("players ON playersId = id").
			Where("externalId = ?", summ.SummonerId()).
			RunWith(db).QueryRow().Scan(&id, &numGames, &goldPerMinute, &csPerMinute,
			&damagePerMinute, &kills, &deaths, &assists, &visionWards, &controlWards)
		if err != nil {
			return err
		}

		summStats := match.PlayerStats(summ)

		gpm := float64(summStats.GoldEarned) / float64(match.GameDuration()*60)
		goldPerMinute = (goldPerMinute*float64(numGames) + float64(gpm)) /
			(float64(numGames + 1))

		//cspm := summStats. / (match.GameDuration() * 60)
		//goldPerMinute = (goldPerMinute * float64(numGames) + float64(gpm)) /
		//	(float64(numGames + 1))

		dpm := float64(summStats.TotalDamageDealtToChampions) / float64(match.GameDuration()*60)
		damagePerMinute = (damagePerMinute*float64(numGames) + dpm) /
			(float64(numGames + 1))

		kills = (kills*float64(numGames) + float64(summStats.Kills)) /
			float64(numGames+1)

		deaths = (deaths*float64(numGames) + float64(summStats.Deaths)) /
			float64(numGames+1)

		assists = (assists*float64(numGames) + float64(summStats.Assists)) /
			float64(numGames+1)

		//visionWards = (visionWards * float64(numGames) + float64(summStats.)) /
		//	float64(numGames + 1)
	}
	return nil
}
