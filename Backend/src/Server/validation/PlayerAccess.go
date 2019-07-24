package validation

import (
	"Server/dataModel"
	"errors"
)

func (a *AccessChecker) Player(accessType AccessType, permissions *dataModel.UserWithPermissions,
	teamDao dataModel.TeamDAO, leagueId, teamId, playerId int) (bool, error) {
	switch accessType {
	case View:
		return teamDao.DoesPlayerExist(leagueId, teamId, playerId)
	case Edit, Delete:
		if permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditTeams ||
			teamAdministrator(permissions, teamId) ||
			teamInformation(permissions, teamId) {
			return teamDao.DoesPlayerExist(leagueId, teamId, playerId)
		} else {
			return false, nil
		}
	case Create:
		return a.Team(Edit, permissions, teamDao, leagueId, teamId)
	default:
		return false, errors.New("invalid access type to check")
	}
}
