package validation

import (
	"Server/dataModel"
	"errors"
)

func (a *AccessChecker) League(accessType AccessType, permissions *dataModel.UserWithPermissions,
	leagueDao dataModel.LeagueDAO, leagueId int) (bool, error) {
	switch accessType {
	case View:
		return leagueDao.DoesLeagueExist(leagueId)
	case Edit:
		return permissions.LeaguePermissions.Administrator, nil
	case Delete:
		return permissions.LeaguePermissions.Administrator, nil
	case Create:
		return true, nil
	default:
		return false, errors.New("invalid access type to check")
	}
}
