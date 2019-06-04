package databaseAccess

const (
	MinGameIdentifierLength = 1
)

func playerGameIdentifierUniquenessValid(leagueId, playerId int, gameIdentifier string, problem *string, errorDest *error) bool {
	var count int
	if err := psql.Select("count(*)").
		From("player").
		Where("league_id = ? AND player_id != ? AND game_identifier = ?", leagueId, playerId, gameIdentifier).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		*errorDest = err
		return false
	} else if count > 0 {
		*problem = "league name already in use"
		return false
	} else {
		return true
	}
}

func playerExternalIdentifierUniquenessValid(leagueId, playerId int, externalId string, problem *string, errorDest *error) bool {
	var count int
	if err := psql.Select("count(*)").
		From("player").
		Where("league_id = ? AND player_id != ? AND external_id = ?", leagueId, playerId, externalId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		*errorDest = err
		return false
	} else if count > 0 {
		*problem = "external id already in use"
		return false
	} else {
		return true
	}
}

func gameIdentifierStringValid(gameIdentifier string, problem *string) bool {
	valid := false
	if len(gameIdentifier) < MinGameIdentifierLength {
		*problem = "game identifier too short"
	} else if len(gameIdentifier) > MaxNameLength {
		*problem = "game identifier too long"
	} else {
		valid = true
	}
	return valid
}
