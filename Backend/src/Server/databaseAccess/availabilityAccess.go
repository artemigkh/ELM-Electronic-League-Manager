package databaseAccess

import "errors"

func (a *AccessChecker) Availability(accessType AccessType, leagueId, availabilityId, userId int) (bool, error) {
	if accessType == Create && availabilityId > 0 {
		return false, errors.New("can't check create permissions for an existing availability")
	}

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
