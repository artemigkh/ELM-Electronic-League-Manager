package databaseAccess

import "errors"

func (a *AccessChecker) Player(accessType AccessType, leagueId, teamId, playerId, userId int) (bool, error) {
	if accessType == Create && playerId > 0 {
		return false, errors.New("can't check create permissions for an existing player")
	}

	// check if player exists on team in db and team exists in league
	var count int
	if err := psql.Select("count(*)").
		From("team").
		Join("player ON player.team_id = team.team_id").
		Join("league ON team.league_id = team.league_id").
		Where("team.league_id = ? AND player.team_id = ? AND player.player_id = ", leagueId, teamId, playerId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else if count == 0 && accessType != Create {
		return false, nil
	}

	// if player exists in this team and league, it is viewable
	if accessType == View {
		return true, nil
	} else {
		leaguePermissions, teamPermissions, err := getLeagueAndTeamPermissions(leagueId, teamId, userId)
		if err != nil {
			return false, err
		}

		admin := leaguePermissions.Administrator || teamPermissions.Administrator

		// go through all cases to catch unsupported access types
		switch accessType {
		case Edit:
			return admin || leaguePermissions.EditTeams || teamPermissions.Players, nil
		case Create:
			return admin || leaguePermissions.EditTeams || teamPermissions.Players, nil
		case Delete:
			return admin, nil
		default:
			return false, errors.New("invalid access type for checker")
		}
	}
}
