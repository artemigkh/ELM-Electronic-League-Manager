package dataModel

type Player struct {
	PlayerId       int    `json:"playerId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}

type PlayerCore struct {
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}

func (player *PlayerCore) validate(leagueId, teamId, playerId int, teamDao TeamDAO) (bool, string, error) {
	return validate(
		validateName(player.Name),
		validateGameIdentifier(player.GameIdentifier),
		player.uniquenessWithExisting(leagueId, playerId, teamDao))
}

func (player *PlayerCore) ValidateNew(leagueId, teamId int, teamDao TeamDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, 0, teamDao)
}

func (player *PlayerCore) ValidateEdit(leagueId, teamId, playerId int, teamDao TeamDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, playerId, teamDao)
}

func (player PlayerCore) uniqueness(playerId int, players []*Player) (bool, string) {
	for _, existingPlayer := range players {
		if existingPlayer.GameIdentifier == player.GameIdentifier && existingPlayer.PlayerId != playerId {
			return false, PlayerGameIdentifierInUse
		}
	}
	return true, ""
}

func (player *PlayerCore) uniquenessWithExisting(leagueId, playerId int, teamDao TeamDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		teams, err := teamDao.GetAllTeamsInLeague(leagueId)
		if err != nil {
			errorDest = &err
			return false
		}

		allPlayers := make([]*Player, 0)
		for _, team := range teams {
			allPlayers = append(allPlayers, team.Players...)
		}

		valid, problem := player.uniqueness(playerId, allPlayers)
		if !valid {
			*problemDest = problem
		}
		return valid
	}
}
