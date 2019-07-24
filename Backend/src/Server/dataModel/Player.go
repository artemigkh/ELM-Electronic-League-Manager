package dataModel

type PlayerCore struct {
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}

func (player *PlayerCore) validate(leagueId, teamId, playerId int, teamDao TeamDAO) (bool, string, error) {
	return validate(
		player.name(),
		player.gameIdentifier(),
		player.uniqueness(leagueId, playerId, teamDao))
}

func (player *PlayerCore) ValidateNew(leagueId, teamId int, teamDao TeamDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, 0, teamDao)
}

func (player *PlayerCore) ValidateEdit(leagueId, teamId, playerId int, teamDao TeamDAO) (bool, string, error) {
	return player.validate(leagueId, teamId, playerId, teamDao)
}

func (player *PlayerCore) name() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(player.Name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(player.Name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (player *PlayerCore) gameIdentifier() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(player.GameIdentifier) > MaxNameLength {
			*problemDest = GameIdentifierTooLong
		} else if len(player.GameIdentifier) < MinInformationLength {
			*problemDest = GameIdentifierTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (player *PlayerCore) uniqueness(leagueId, playerId int, teamDao TeamDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		teams, err := teamDao.GetAllTeamsInLeague(leagueId)
		if err != nil {
			errorDest = &err
			return false
		}

		for _, team := range teams {
			for _, existingPlayer := range team.Players {
				if existingPlayer.GameIdentifier == player.GameIdentifier && existingPlayer.PlayerId != playerId {
					*problemDest = PlayerGameIdentifierInUse
					return false
				}
			}
		}
		return true
	}
}

type Player struct {
	PlayerId       int    `json:"playerId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}
