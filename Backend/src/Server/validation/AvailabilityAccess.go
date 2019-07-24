package validation

import (
	"Server/dataModel"
	"errors"
)

func (a *AccessChecker) Availability(accessType AccessType, permissions *dataModel.UserWithPermissions,
	leagueDao dataModel.LeagueDAO, leagueId, availabilityId int) (bool, error) {
	switch accessType {
	case View:
		return leagueDao.DoesAvailabilityExistInLeague(leagueId, availabilityId)
	case Edit, Delete:
		if permissions.LeaguePermissions.Administrator || permissions.LeaguePermissions.EditGames {
			return leagueDao.DoesAvailabilityExistInLeague(leagueId, availabilityId)
		} else {
			return false, nil
		}
	case Create:
		return permissions.LeaguePermissions.Administrator || permissions.LeaguePermissions.EditGames, nil
	default:
		return false, errors.New("invalid access type to check")
	}
}
