package validation

import (
	"Server/dataModel"
	"errors"
)

func (a *AccessChecker) Game(accessType AccessType, permissions *dataModel.UserWithPermissions,
	gameDao dataModel.GameDAO, leagueId, gameId int) (bool, error) {
	switch accessType {
	case View:
		return gameDao.DoesGameExistInLeague(leagueId, gameId)
	case Edit:
		gameExists, err := gameDao.DoesGameExistInLeague(leagueId, gameId)
		if err != nil {
			return false, err
		} else if !gameExists {
			return false, nil
		}
		if permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditGames {
			return true, nil
		} else {
			game, err := gameDao.GetGameInformation(gameId)
			if err != nil {
				return false, err
			}
			return teamAdministrator(permissions, game.Team1.TeamId) ||
				teamAdministrator(permissions, game.Team2.TeamId) ||
				teamGames(permissions, game.Team1.TeamId) ||
				teamGames(permissions, game.Team2.TeamId), nil
		}
	case Delete:
		if permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditGames {
			return gameDao.DoesGameExistInLeague(leagueId, gameId)
		} else {
			return false, nil
		}
	case Create:
		return permissions.LeaguePermissions.Administrator ||
			permissions.LeaguePermissions.EditGames, nil
	default:
		return false, errors.New("invalid access type to check")
	}
}
