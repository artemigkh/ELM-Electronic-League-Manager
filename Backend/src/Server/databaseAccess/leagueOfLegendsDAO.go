package databaseAccess

import (
	"Server/lolApi"
	"database/sql"
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

func (d *PgLeagueOfLegendsDAO) createChampionStatsIfNotExist(leagueId int, champion string) error {
	// check if exists
	var id int
	err := psql.Select("leagueId").
		From("championStats").
		Where("leagueId = ? AND name = ?", leagueId, champion).
		RunWith(db).QueryRow().
		Scan(&id)
	if err == sql.ErrNoRows {
		// does not exist, so create
		_, err := psql.Insert("championStats").
			Columns("leagueId", "name", "picks", "wins", "bans").
			Values(leagueId, champion, 0, 0, 0).RunWith(db).Exec()
		if err != nil {
			return err
		}
	} else if err != nil {
		// some db error occured
		return err
	}

	//exists so do nothing
	return nil
}

func (d *PgLeagueOfLegendsDAO) updateChampionStats(leagueId int, match *lolApi.MatchInformation) error {
	for _, champion := range match.BannedChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			`
		UPDATE championStats SET bans = bans + 1
		WHERE leagueId = $1 AND name = $2
		`, leagueId, champion)
		if err != nil {
			return err
		}
	}

	for _, champion := range match.WinningChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			`
		UPDATE championStats SET picks = picks + 1, wins = wins + 1
		WHERE leagueId = $1 AND name = $2
		`, leagueId, champion)
		if err != nil {
			return err
		}
	}

	for _, champion := range match.LosingChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			`
		UPDATE championStats SET picks = picks + 1
		WHERE leagueId = $1 AND name = $2
		`, leagueId, champion)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *PgLeagueOfLegendsDAO) ReportEndGameStats(leagueId, gameId,
	winTeamId, loseTeamId int, match *lolApi.MatchInformation) error {

	err := d.updateChampionStats(leagueId, match)
	if err != nil {
		return err
	}

	// Create League Game Entry
	var leagueGameId int
	err = psql.Insert("leagueGame").
		Columns("gameId", "leagueId", "winTeamId", "loseTeamId", "timestamp", "duration").
		Values(gameId, leagueId, winTeamId, loseTeamId, match.Timestamp, match.Duration).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueGameId)
	if err != nil {
		return err
	}

	// Create Winning Team Stats Entry for this game
	_, err = psql.Insert("teamStats").
		Columns("teamId", "gameId", "leagueId", "duration", "side", "firstBlood", "firstTurret", "win").
		Values(winTeamId, leagueGameId, leagueId, match.Duration, match.WinningTeamStats.Side,
			match.WinningTeamStats.FirstBlood, match.WinningTeamStats.FirstTower, true).RunWith(db).Exec()
	if err != nil {
		return err
	}

	// Create Losing Team Stats Entry for this game
	_, err = psql.Insert("teamStats").
		Columns("teamId", "gameId", "leagueId", "duration", "side", "firstBlood", "firstTurret", "win").
		Values(loseTeamId, leagueGameId, leagueId, match.Duration, match.LosingTeamStats.Side,
			match.LosingTeamStats.FirstBlood, match.LosingTeamStats.FirstTower, false).RunWith(db).Exec()
	if err != nil {
		return err
	}

	// Create Stats Entry for each Player
	for _, player := range match.PlayerStats {
		var teamId int
		if player.Win {
			teamId = winTeamId
		} else {
			teamId = loseTeamId
		}
		_, err = psql.Insert("playerStats").
			Columns("id", "name", "gameId", "teamId", "leagueId", "duration", "championPicked",
				"gold", "cs", "damage", "kills", "deaths", "assists", "wards", "win").
			Values(player.Id, player.Name, leagueGameId, teamId, leagueId, match.Duration,
				player.ChampionPicked, player.Gold, player.Cs, player.Damage,
				player.Kills, player.Deaths, player.Assists, player.Wards, player.Win).RunWith(db).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
