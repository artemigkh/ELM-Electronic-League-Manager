package validation

import (
	"Server/dataModel"
	"github.com/pkg/errors"
)

func teamPermissions(userPermissions *dataModel.UserWithPermissions, teamId int) *dataModel.TeamPermissions {
	for _, teamPermissions := range userPermissions.TeamPermissions {
		if teamPermissions.TeamId == teamId {
			return teamPermissions
		}
	}
	return nil
}

func teamAdministrator(userPermissions *dataModel.UserWithPermissions, teamId int) bool {
	if permissions := teamPermissions(userPermissions, teamId); permissions == nil {
		return false
	} else {
		return permissions.Administrator
	}
}

func teamInformation(userPermissions *dataModel.UserWithPermissions, teamId int) bool {
	if permissions := teamPermissions(userPermissions, teamId); permissions == nil {
		return false
	} else {
		return permissions.Information
	}
}

func teamGames(userPermissions *dataModel.UserWithPermissions, teamId int) bool {
	if permissions := teamPermissions(userPermissions, teamId); permissions == nil {
		return false
	} else {
		return permissions.Games
	}
}

func (a *AccessChecker) Team(accessType AccessType, permissions *dataModel.UserWithPermissions,
	teamDao dataModel.TeamDAO, leagueId, teamId int) (bool, error) {
	switch accessType {
	case View:
		return teamDao.DoesTeamExistInLeague(leagueId, teamId)
	case Edit:
		if teamAdministrator(permissions, teamId) ||
			teamInformation(permissions, teamId) {
			return true, nil
		} else if permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditTeams {
			return teamDao.DoesTeamExistInLeague(leagueId, teamId)
		} else {
			return false, nil
		}
	case Delete:
		if teamAdministrator(permissions, teamId) {
			return true, nil
		} else if permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditTeams {
			return teamDao.DoesTeamExistInLeague(leagueId, teamId)
		} else {
			return false, nil
		}
	case Create:
		return permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.CreateTeams, nil
	default:
		return false, errors.New("invalid access type to check")
	}
}
