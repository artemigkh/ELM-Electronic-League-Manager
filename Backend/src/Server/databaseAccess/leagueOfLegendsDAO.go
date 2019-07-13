package databaseAccess

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
	AverageWards    float64 `json:"averageWards"`
}

type TeamStats struct {
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

type ChampionStats struct {
	Name    string  `json:"name"`
	Bans    int     `json:"bans"`
	Picks   int     `json:"picks"`
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	Winrate float64 `json:"winrate"`
}

type PgLeagueOfLegendsDAO struct{}

func (d *PgLeagueOfLegendsDAO) createChampionStatsIfNotExist(leagueId int, champion string) error {
	//// check if exists
	//var id int
	//err := psql.Select("league_id").
	//	From("championStats").
	//	Where("league_id = ? AND name = ?", leagueId, champion).
	//	RunWith(db).QueryRow().
	//	Scan(&id)
	//if err == sql.ErrNoRows {
	//	// does not exist, so create
	//	_, err := psql.Insert("championStats").
	//		Columns("league_id", "name", "picks", "wins", "bans").
	//		Values(leagueId, champion, 0, 0, 0).RunWith(db).Exec()
	//	if err != nil {
	//		return err
	//	}
	//} else if err != nil {
	//	// some db error occured
	//	return err
	//}
	//
	////exists so do nothing
	return nil
}

//func (d *PgLeagueOfLegendsDAO) updateChampionStats(leagueId int, match *lolApi.MatchInformation) error {
//for _, champion := range match.BannedChampions {
//	err := d.createChampionStatsIfNotExist(leagueId, champion)
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(
//		`
//	UPDATE championStats SET bans = bans + 1
//	WHERE league_id = $1 AND name = $2
//	`, leagueId, champion)
//	if err != nil {
//		return err
//	}
//}
//
//for _, champion := range match.WinningChampions {
//	err := d.createChampionStatsIfNotExist(leagueId, champion)
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(
//		`
//	UPDATE championStats SET picks = picks + 1, wins = wins + 1
//	WHERE league_id = $1 AND name = $2
//	`, leagueId, champion)
//	if err != nil {
//		return err
//	}
//}
//
//for _, champion := range match.LosingChampions {
//	err := d.createChampionStatsIfNotExist(leagueId, champion)
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(
//		`
//	UPDATE championStats SET picks = picks + 1
//	WHERE league_id = $1 AND name = $2
//	`, leagueId, champion)
//	if err != nil {
//		return err
//	}
//}

//	return nil
//}

//func (d *PgLeagueOfLegendsDAO) ReportEndGameStats(leagueId, gameId,
//	winTeamId, loseTeamId int, match *lolApi.MatchInformation) error {
//	return nil
//
//err := d.updateChampionStats(leagueId, match)
//if err != nil {
//	return err
//}
//
//// Create League Game Entry
//var leagueGameId int
//err = psql.Insert("leagueGame").
//	Columns("gameId", "league_id", "winTeamId", "loseTeamId", "timestamp", "duration").
//	Values(gameId, leagueId, winTeamId, loseTeamId, match.Timestamp, match.Duration).
//	Suffix("RETURNING \"id\"").
//	RunWith(db).QueryRow().Scan(&leagueGameId)
//if err != nil {
//	return err
//}
//
//// Create Winning Team Stats Entry for this game
//_, err = psql.Insert("teamStats").
//	Columns("team_id", "gameId", "league_id", "duration", "side", "firstBlood", "firstTurret", "win").
//	Values(winTeamId, leagueGameId, leagueId, match.Duration, match.WinningTeamStats.Side,
//		match.WinningTeamStats.FirstBlood, match.WinningTeamStats.FirstTower, true).RunWith(db).Exec()
//if err != nil {
//	return err
//}
//
//// Create Losing Team Stats Entry for this game
//_, err = psql.Insert("teamStats").
//	Columns("team_id", "gameId", "league_id", "duration", "side", "firstBlood", "firstTurret", "win").
//	Values(loseTeamId, leagueGameId, leagueId, match.Duration, match.LosingTeamStats.Side,
//		match.LosingTeamStats.FirstBlood, match.LosingTeamStats.FirstTower, false).RunWith(db).Exec()
//if err != nil {
//	return err
//}
//
//// Create Stats Entry for each Player
//for _, player := range match.PlayerStats {
//	var teamId int
//	if player.Win {
//		teamId = winTeamId
//	} else {
//		teamId = loseTeamId
//	}
//	_, err = psql.Insert("playerStats").
//		Columns("id", "name", "gameId", "team_id", "league_id", "duration", "championPicked",
//			"gold", "cs", "damage", "kills", "deaths", "assists", "wards", "win").
//		Values(player.Id, player.Name, leagueGameId, teamId, leagueId, match.Duration,
//			player.ChampionPicked, player.Gold, player.Cs, player.Damage,
//			player.Kills, player.Deaths, player.Assists, player.Wards, player.Win).RunWith(db).Exec()
//	if err != nil {
//		return err
//	}
//}
//
//return nil
//}

func (d *PgLeagueOfLegendsDAO) GetPlayerStats(leagueId int) ([]*PlayerStats, error) {
	return nil, nil
	//	rows, err := db.Query(`
	//SELECT id, (array_agg(name ORDER BY name))[1] as name, (array_agg(team_id ORDER BY team_id))[1] as team_id,
	//SUM(damage) / (SUM(duration) / 60) AS DPM,
	//SUM(gold) / (SUM(duration) / 60) AS GPM,
	//SUM(cs) / (SUM(duration) / 60) AS CSPM,
	//AVG(duration) AS AverageDuration,
	//AVG(kills) AS AverageKills,
	//AVG(deaths) AS AverageDeaths,
	//AVG(assists) AS AverageAssists,
	//AVG(wards) AS AverageWards,
	//(AVG(kills) + AVG(assists)) / GREATEST(1, AVG(deaths)) AS AverageKda
	//FROM playerStats WHERE league_id = $1
	//GROUP BY id`, leagueId)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	defer rows.Close()
	//
	//	var allPlayerStats []*PlayerStats
	//
	//	for rows.Next() {
	//		var playerStats PlayerStats
	//		err := rows.Scan(&playerStats.Id, &playerStats.Name, &playerStats.TeamId, &playerStats.DamagePerMinute,
	//			&playerStats.GoldPerMinute, &playerStats.CsPerMinute, &playerStats.AverageDuration,
	//			&playerStats.AverageKills, &playerStats.AverageDeaths, &playerStats.AverageAssists,
	//			&playerStats.AverageWards, &playerStats.AverageKda)
	//		if err != nil {
	//			return nil, err
	//		}
	//		allPlayerStats = append(allPlayerStats, &playerStats)
	//	}
	//	err = rows.Err()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return allPlayerStats, nil
	//}

	//type TeamStats struct {
	//	Id                 string  `json:"id"`
	//	AverageDuration    float64 `json:"averageDuration"`
	//	NumberFirstBloods  int     `json:"numberFirstBloods"`
	//	NumberFirstTurrets int     `json:"numberFirstTurrets"`
	//}
	//func (d *PgLeagueOfLegendsDAO) GetTeamStats(leagueId int) ([]*TeamStats, error) {
	//	return nil, nil
	//	rows, err := db.Query(`
	//SELECT t1.team_id, averageDuration, numberFirstBloods, numberFirstTurrets,
	//AverageKda, AverageActionScore, AverageWards, AverageGoldPerMinute, AverageCsPerMinute
	//FROM (SELECT team_id,
	//(AVG(kills) + AVG(assists)) / GREATEST(1, AVG(deaths)) AS AverageKda,
	//(SUM(kills) + SUM(deaths)) / (COUNT(*)/5) AS AverageActionScore,
	//SUM(wards) / (COUNT(*)/5) AS AverageWards,
	//(SUM(gold) / (SUM(duration)/60)) / (COUNT(*)/5) AS AverageGoldPerMinute,
	//(SUM(cs) / (SUM(duration)/60)) / (COUNT(*)/5) AS AverageCsPerMinute
	//FROM playerStats WHERE league_id=$1
	//GROUP BY team_id) AS t1
	//INNER JOIN
	//(SELECT team_id, AVG(duration) AS averageDuration,
	//COUNT(*) FILTER (WHERE firstBlood) AS numberFirstBloods,
	//COUNT(*) FILTER (WHERE firstTurret) AS numberFirstTurrets
	//FROM teamStats WHERE league_id = $1
	//GROUP BY team_id) AS t2
	//ON t1.team_id = t2.team_id`, leagueId)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	defer rows.Close()
	//
	//	var allTeamStats []*TeamStats
	//	for rows.Next() {
	//		var teamStats TeamStats
	//		err := rows.Scan(&teamStats.Id, &teamStats.AverageDuration, &teamStats.NumberFirstBloods,
	//			&teamStats.NumberFirstTurrets, &teamStats.AverageKda, &teamStats.AverageActionScore,
	//			&teamStats.AverageWards, &teamStats.GoldPerMinute, &teamStats.CsPerMinute)
	//		if err != nil {
	//			return nil, err
	//		}
	//		allTeamStats = append(allTeamStats, &teamStats)
	//	}
	//	err = rows.Err()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return allTeamStats, nil
}

//
//func (d *PgLeagueOfLegendsDAO) GetChampionStats(leagueId int) ([]*ChampionStats, error) {
//	return nil, nil
//	rows, err := db.Query(`
//SELECT name, bans, picks, wins, picks-wins AS losses,
//CASE picks
//  WHEN 0 THEN 0
//  ELSE wins::FLOAT / picks::FLOAT
//END AS winrate
//FROM championStats WHERE league_id = $1`, leagueId)
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	var allChampStats []*ChampionStats
//
//	for rows.Next() {
//		var champStats ChampionStats
//		err := rows.Scan(&champStats.Name, &champStats.Bans, &champStats.Picks, &champStats.Wins,
//			&champStats.Losses, &champStats.Winrate)
//		if err != nil {
//			return nil, err
//		}
//		allChampStats = append(allChampStats, &champStats)
//	}
//	err = rows.Err()
//	if err != nil {
//		return nil, err
//	}
//
//	return allChampStats, nil
//}

func (d *PgLeagueOfLegendsDAO) CreateLoLPlayer(leagueId, teamId int, externalId string, playerInfo LoLPlayerCore) (int, error) {
	var playerId int
	if err := psql.Insert("player").
		Columns(
			"team_id",
			"league_id",
			"game_identifier",
			"name",
			"main_roster",
			"external_id",
			"position",
		).
		Values(
			teamId,
			leagueId,
			playerInfo.GameIdentifier,
			playerInfo.GameIdentifier,
			playerInfo.MainRoster,
			externalId,
			playerInfo.Position,
		).
		Suffix("RETURNING \"player_id\"").
		RunWith(db).QueryRow().Scan(&playerId); err != nil {
		return -1, err
	}

	return playerId, nil
}

func (d *PgLeagueOfLegendsDAO) UpdateLoLPlayer(playerId int, externalId string, playerInfo LoLPlayerCore) error {
	_, err := psql.Update("player").
		Set("game_identifier", playerInfo.GameIdentifier).
		Set("name", playerInfo.GameIdentifier).
		Set("main_roster", playerInfo.MainRoster).
		Set("external_id", externalId).
		Set("position", playerInfo.Position).
		Where("player_id = ?", playerId).
		RunWith(db).Exec()

	return err
}

func (d *PgLeagueOfLegendsDAO) GetLoLTeamStub(teamId int) (*LoLTeamStub, error) {
	rows, err := getLoLTeamStubSelector().
		Where("team.team_id = ?", teamId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedLoLTeamStub(rows)
}
