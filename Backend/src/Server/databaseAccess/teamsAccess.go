package databaseAccess

import "errors"

func (a *AccessChecker) Team(accessType AccessType, leagueId, teamId, userId int) (bool, error) {
	if accessType == Create && teamId > 0 {
		return false, errors.New("can't check create permissions for an existing team")
	}
	// check if team exists in league
	var count int
	if err := psql.Select("count(*)").
		From("team").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else if count == 0 && accessType != Create {
		return false, nil
	}
	// if team exists in this league, it is viewable
	if accessType == View {
		return true, nil
	} else {
		leaguePermissions, teamPermissions, err := getLeagueAndTeamPermissions(leagueId, teamId, userId)
		if err != nil {
			return false, err
		}

		admin := leaguePermissions.Administrator || teamPermissions.Administrator
		if admin {
			return true, nil
		}

		// go through all cases to catch unsupported access types
		switch accessType {
		case Edit:
			return admin || leaguePermissions.EditTeams || teamPermissions.Information, nil
		case Create:
			return admin || leaguePermissions.CreateTeams, nil
		case Delete:
			return admin, nil
		default:
			return false, errors.New("invalid access type for checker")
		}
	}
}
