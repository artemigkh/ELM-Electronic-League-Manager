package databaseAccess

import "errors"

func (a *AccessChecker) Game(accessType AccessType, leagueId, gameId, userId int) (bool, error) {
	if accessType == Create && gameId > 0 {
		return false, errors.New("can't check create permissions for an existing game")
	}

	// check if game exists in league
	var count int
	if err := psql.Select("count(*)").
		From("game").
		Where("league_id = ? AND game_id = ?", leagueId, gameId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else if count == 0 && accessType != Create {
		return false, nil
	}

	// if game exists in this league, it is viewable
	if accessType == View {
		return true, nil
	} else {
		leaguePermissions, err := getLeaguePermissions(leagueId, userId)
		if err != nil {
			return false, err
		}

		return leaguePermissions.Administrator || leaguePermissions.EditGames, nil
	}
}

func (a *AccessChecker) Report(leagueId, gameId, userId int) (bool, error) {
	return false, nil
}
