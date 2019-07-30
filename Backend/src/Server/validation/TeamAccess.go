package validation

import (
	"Server/dataModel"
	"github.com/pkg/errors"
	"time"
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
	teamDao dataModel.TeamDAO, leagueDao dataModel.LeagueDAO, leagueId, teamId int) (bool, error) {
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
		if permissions.LeaguePermissions.Administrator || permissions.LeaguePermissions.CreateTeams {
			return true, nil
		} else {
			// if during signup period and allows public signups
			leagueInfo, err := leagueDao.GetLeagueInformation(leagueId)
			if err != nil {
				return false, err
			} else {
				return leagueInfo.PublicJoin &&
					time.Now().Unix() > int64(leagueInfo.SignupStart) &&
					time.Now().Unix() < int64(leagueInfo.SignupEnd), nil
			}
		}
	default:
		return false, errors.New("invalid access type to check")
	}
}
