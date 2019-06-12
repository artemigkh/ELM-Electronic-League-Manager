package databaseAccess

import "github.com/pkg/errors"

func (a *AccessChecker) League(accessType AccessType, leagueId, userId int) (bool, error) {
	if accessType == Create && leagueId > 0 {
		return false, errors.New("can't check create permissions for an existing league")
	}
	// check if team league exists
	var count int
	if err := psql.Select("count(*)").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else if count == 0 && accessType != Create {
		return false, nil
	}

	if accessType == Create || accessType == View {
		return true, nil
	} else {
		leaguePermissions, err := getLeaguePermissions(leagueId, userId)
		if err != nil {
			return false, err
		}

		switch accessType {
		case Edit:
			return leaguePermissions.Administrator, nil
		case Delete:
			return false, nil
		default:
			return false, errors.New("invalid access type for checker")
		}
	}
}
