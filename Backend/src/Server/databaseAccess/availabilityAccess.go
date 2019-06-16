package databaseAccess

import "errors"

func (a *AccessChecker) Availability(accessType AccessType, leagueId, availabilityId, userId int) (bool, error) {
	if accessType == Create && availabilityId > 0 {
		return false, errors.New("can't check create permissions for an existing availability")
	}

	// check if availability exists in league
	var count int
	if err := psql.Select("count(*)").
		From("availability").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else if count == 0 && accessType != Create {
		return false, nil
	}

	// if availability exists in this league, it is viewable
	if accessType == View {
		return true, nil
	} else {
		leaguePermissions, err := getLeaguePermissions(leagueId, userId)
		if err != nil {
			return false, err
		}

		switch accessType {
		case Create:
			return leaguePermissions.Administrator, nil
		case Edit:
			return leaguePermissions.Administrator, nil
		case Delete:
			return leaguePermissions.Administrator, nil
		default:
			return false, errors.New("invalid access type for checker")
		}
	}
}
