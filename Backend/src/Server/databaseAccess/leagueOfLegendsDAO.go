package databaseAccess

import (
	"Server/lolApi"
	"database/sql"
)

type PlayerStats struct {
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
}

type TeamStats struct {
	Id                 string  `json:"id"`
	AverageDuration    float64 `json:"averageDuration"`
	NumberFirstBloods  int     `json:"numberFirstBloods"`
	NumberFirstTurrets int     `json:"numberFirstTurrets"`
}

type ChampionStats struct {
	Name    string  `json:"name"`
	Picks   int     `json:"picks"`
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	Winrate float64 `json:"winrate"`
}

type PgLeagueOfLegendsDAO struct{}

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

func (d *PgLeagueOfLegendsDAO) GetPlayerStats(leagueId int) ([]*PlayerStats, error) {
	rows, err := db.Query(`
SELECT id, (array_agg(name ORDER BY name))[1] as name, (array_agg(teamId ORDER BY teamId))[1] as teamId,
SUM(damage) / (SUM(duration) / 60) AS DPM,
SUM(gold) / (SUM(duration) / 60) AS GPM,
SUM(cs) / (SUM(duration) / 60) AS CSPM,
AVG(duration) AS AverageDuration,
AVG(kills) AS AverageKills,
AVG(deaths) AS AverageDeaths,
AVG(assists) AS AverageAssists,	
(AVG(kills) + AVG(assists)) / GREATEST(1, AVG(deaths)) AS AverageKda
FROM playerStats WHERE leagueId = $1
GROUP BY id`, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allPlayerStats []*PlayerStats
	var playerStats PlayerStats

	for rows.Next() {
		err := rows.Scan(&playerStats.Id, &playerStats.Name, &playerStats.TeamId, &playerStats.DamagePerMinute,
			&playerStats.GoldPerMinute, &playerStats.CsPerMinute, &playerStats.AverageDuration,
			&playerStats.AverageKills, &playerStats.AverageDeaths, &playerStats.AverageAssists, &playerStats.AverageKda)
		if err != nil {
			return nil, err
		}
		allPlayerStats = append(allPlayerStats, &playerStats)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allPlayerStats, nil
}

func (d *PgLeagueOfLegendsDAO) GetTeamStats(leagueId int) ([]*TeamStats, error) {
	return nil, nil
}

func (d *PgLeagueOfLegendsDAO) GetChampionStats(leagueId int) ([]*ChampionStats, error) {
	return nil, nil
}
