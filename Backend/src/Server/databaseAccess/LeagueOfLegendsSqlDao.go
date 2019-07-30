package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type LeagueOfLegendsSqlDao struct{}

func (d *LeagueOfLegendsSqlDao) createChampionStatsIfNotExist(leagueId int, champion string) error {
	// check if exists
	var id int
	err := psql.Select("league_id").
		From("lol_champion_stats").
		Where("league_id = ? AND name = ?", leagueId, champion).
		RunWith(db).QueryRow().
		Scan(&id)
	if err == sql.ErrNoRows {
		// does not exist, so create
		_, err := psql.Insert("lol_champion_stats").
			Columns("league_id", "name", "picks", "wins", "bans").
			Values(leagueId, champion, 0, 0, 0).RunWith(db).Exec()
		if err != nil {
			return err
		}
	} else if err != nil {
		// some db error occurred
		return err
	}

	//exists so do nothing
	return nil
}

func (d *LeagueOfLegendsSqlDao) updateChampionStats(leagueId int, match *dataModel.LoLMatchInformation) error {
	for _, champion := range match.BannedChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = psql.Update("lol_champion_stats").
			Set("bans", squirrel.Expr("bans + 1")).
			Where("league_id = ? AND name = ?", leagueId, champion).
			RunWith(db).Exec()
		if err != nil {
			return err
		}
	}

	for _, champion := range match.WinningChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = psql.Update("lol_champion_stats").
			Set("picks", squirrel.Expr("picks + 1")).
			Set("wins", squirrel.Expr("wins + 1")).
			Where("league_id = ? AND name = ?", leagueId, champion).
			RunWith(db).Exec()
		if err != nil {
			return err
		}
	}

	for _, champion := range match.LosingChampions {
		err := d.createChampionStatsIfNotExist(leagueId, champion)
		if err != nil {
			return err
		}

		_, err = psql.Update("lol_champion_stats").
			Set("picks", squirrel.Expr("picks + 1")).
			Set("wins", squirrel.Expr("wins + 1")).
			Where("league_id = ? AND name = ?", leagueId, champion).
			RunWith(db).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *LeagueOfLegendsSqlDao) ReportEndGameStats(leagueId, gameId int, match *dataModel.LoLMatchInformation) error {
	if err := d.updateChampionStats(leagueId, match); err != nil {
		return err
	}

	// Create Winning Team Stats Entry for this game
	_, err := psql.Insert("lol_team_stats").
		Columns(
			"team_id",
			"game_id",
			"league_id",
			"duration",
			"side",
			"first_blood",
			"first_turret",
			"win",
		).
		Values(
			match.WinningTeamId,
			gameId,
			leagueId,
			match.Duration,
			match.WinningTeamStats.Side,
			match.WinningTeamStats.FirstBlood,
			match.WinningTeamStats.FirstTower,
			true,
		).RunWith(db).Exec()
	if err != nil {
		return err
	}

	// Create Losing Team Stats Entry for this game
	_, err = psql.Insert("lol_team_stats").
		Columns(
			"team_id",
			"game_id",
			"league_id",
			"duration",
			"side",
			"first_blood",
			"first_turret",
			"win",
		).
		Values(
			match.LosingTeamId,
			gameId,
			leagueId,
			match.Duration,
			match.LosingTeamStats.Side,
			match.LosingTeamStats.FirstBlood,
			match.LosingTeamStats.FirstTower,
			false,
		).RunWith(db).Exec()
	if err != nil {
		return err
	}

	// Create Stats Entry for each Player
	for _, player := range match.PlayerStats {
		var teamId int
		if player.Win {
			teamId = match.WinningTeamId
		} else {
			teamId = match.LosingTeamId
		}
		_, err = psql.Insert("lol_player_stats").
			Columns(
				"id",
				"name",
				"game_id",
				"team_id",
				"league_id",
				"duration",
				"champion_picked",
				"gold",
				"cs",
				"damage",
				"kills",
				"deaths",
				"assists",
				"wards",
				"win",
			).
			Values(
				player.Id,
				player.Name,
				gameId,
				teamId,
				leagueId,
				match.Duration,
				player.ChampionPicked,
				player.Gold,
				player.Cs,
				player.Damage,
				player.Kills,
				player.Deaths,
				player.Assists,
				player.Wards,
				player.Win,
			).RunWith(db).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *LeagueOfLegendsSqlDao) GetPlayerStats(leagueId int) ([]*dataModel.LoLPlayerStats, error) {
	rows, err := db.Query(`
	SELECT id, (array_agg(name ORDER BY name))[1] as name, (array_agg(team_id ORDER BY team_id))[1] as team_id,
	SUM(damage) / (SUM(duration) / 60) AS DPM,
	SUM(gold) / (SUM(duration) / 60) AS GPM,
	SUM(cs) / (SUM(duration) / 60) AS CSPM,
	AVG(duration) AS average_duration,
	AVG(kills) AS average_kills,
	AVG(deaths) AS average_deaths,
	AVG(assists) AS average_assists,
	AVG(wards) AS average_wards,
	(AVG(kills) + AVG(assists)) / GREATEST(1, AVG(deaths)) AS average_kda
	FROM lol_player_stats WHERE league_id = $1
	GROUP BY id`, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allPlayerStats := make([]*dataModel.LoLPlayerStats, 0)

	for rows.Next() {
		var playerStats dataModel.LoLPlayerStats
		if err := rows.Scan(
			&playerStats.Id,
			&playerStats.Name,
			&playerStats.TeamId,
			&playerStats.DamagePerMinute,
			&playerStats.GoldPerMinute,
			&playerStats.CsPerMinute,
			&playerStats.AverageDuration,
			&playerStats.AverageKills,
			&playerStats.AverageDeaths,
			&playerStats.AverageAssists,
			&playerStats.AverageWards,
			&playerStats.AverageKda,
		); err != nil {
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

func (d *LeagueOfLegendsSqlDao) GetTeamStats(leagueId int) ([]*dataModel.LoLTeamStats, error) {
	rows, err := db.Query(`
	SELECT t1.team_id, average_duration, number_first_bloods, number_first_turrets,
	average_kda, average_action_score, average_wards, average_gold_per_minute, average_cs_per_minute
	FROM (SELECT team_id,
	(AVG(kills) + AVG(assists)) / GREATEST(1, AVG(deaths)) AS average_kda,
	(SUM(kills) + SUM(deaths)) / (COUNT(*)/5) AS average_action_score,
	SUM(wards) / (COUNT(*)/5) AS average_wards,
	(SUM(gold) / (SUM(duration)/60)) / (COUNT(*)/5) AS average_gold_per_minute,
	(SUM(cs) / (SUM(duration)/60)) / (COUNT(*)/5) AS average_cs_per_minute
	FROM lol_player_stats WHERE league_id=$1
	GROUP BY team_id) AS t1
	INNER JOIN
	(SELECT team_id, AVG(duration) AS average_duration,
	COUNT(*) FILTER (WHERE first_blood) AS number_first_bloods,
	COUNT(*) FILTER (WHERE first_turret) AS number_first_turrets
	FROM lol_team_stats WHERE league_id = $1
	GROUP BY team_id) AS t2
	ON t1.team_id = t2.team_id`, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allTeamStats := make([]*dataModel.LoLTeamStats, 0)
	for rows.Next() {
		var teamStats dataModel.LoLTeamStats
		err := rows.Scan(
			&teamStats.Id,
			&teamStats.AverageDuration,
			&teamStats.NumberFirstBloods,
			&teamStats.NumberFirstTurrets,
			&teamStats.AverageKda,
			&teamStats.AverageActionScore,
			&teamStats.AverageWards,
			&teamStats.GoldPerMinute,
			&teamStats.CsPerMinute,
		)
		if err != nil {
			return nil, err
		}
		allTeamStats = append(allTeamStats, &teamStats)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allTeamStats, nil
}

func (d *LeagueOfLegendsSqlDao) GetChampionStats(leagueId int) ([]*dataModel.LoLChampionStats, error) {
	rows, err := db.Query(`
	SELECT name, bans, picks, wins, picks-wins AS losses,
	CASE picks
	 WHEN 0 THEN 0
	 ELSE wins::FLOAT / picks::FLOAT
	END AS winrate
	FROM lol_champion_stats WHERE league_id = $1`, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allChampStats := make([]*dataModel.LoLChampionStats, 0)

	for rows.Next() {
		var champStats dataModel.LoLChampionStats
		err := rows.Scan(
			&champStats.Name,
			&champStats.Bans,
			&champStats.Picks,
			&champStats.Wins,
			&champStats.Losses,
			&champStats.Winrate,
		)
		if err != nil {
			return nil, err
		}
		allChampStats = append(allChampStats, &champStats)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allChampStats, nil
}

func (d *LeagueOfLegendsSqlDao) CreateLoLPlayer(leagueId, teamId int, externalId string, playerInfo dataModel.LoLPlayerCore) (int, error) {
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

func (d *LeagueOfLegendsSqlDao) UpdateLoLPlayer(playerId int, externalId string, playerInfo dataModel.LoLPlayerCore) error {
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

func (d *LeagueOfLegendsSqlDao) GetLoLTeamStub(teamId int) (*dataModel.LoLTeamStub, error) {
	rows, err := getLoLTeamStubSelector().
		Where("team.team_id = ?", teamId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedLoLTeamStub(rows)
}

func (d *LeagueOfLegendsSqlDao) GetAllLoLTeamStubInLeague(leagueId int) ([]*dataModel.LoLTeamStub, error) {
	rows, err := getLoLTeamStubSelector().
		Where("team.league_id = ?", leagueId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedAllLoLTeamStubs(rows)
}

func (d *LeagueOfLegendsSqlDao) CreateLoLTeamWithPlayers(
	leagueId,
	userId int,
	teamInfo dataModel.TeamCore,
	players []*dataModel.LoLPlayerCore,
	iconSmall, iconLarge string) (int, error) {

	// Create new team based on icon status
	var teamId int
	var err error
	if len(iconSmall)+len(iconLarge) > 0 {
		teamId, err = teamsDAO.CreateTeamWithIcon(leagueId, userId, teamInfo, iconSmall, iconLarge)
	} else {
		teamId, err = teamsDAO.CreateTeam(leagueId, userId, teamInfo)
	}
	if err != nil {
		return -1, err
	}

	if players == nil || len(players) == 0 {
		return teamId, nil
	}
	// Insert all players
	insertBuilder := psql.Insert("player").
		Columns(
			"team_id",
			"league_id",
			"game_identifier",
			"name",
			"main_roster",
			"external_id",
			"position",
		)

	for _, player := range players {
		insertBuilder = insertBuilder.Values(
			teamId,
			leagueId,
			player.GameIdentifier,
			player.GameIdentifier,
			player.MainRoster,
			player.ExternalId,
			player.Position,
		)
	}

	_, err = insertBuilder.RunWith(db).Exec()
	if err != nil {
		//TODO: make sure permission table entries also deleted here with cascade

		//"rollback" delete the created team if this insertion failed
		_, delErr := psql.Delete("team").
			Where("team_id = ?", teamId).
			RunWith(db).Exec()
		if delErr != nil {
			return -1, delErr
		} else {
			return -1, err
		}
	} else {
		return teamId, nil
	}
}
