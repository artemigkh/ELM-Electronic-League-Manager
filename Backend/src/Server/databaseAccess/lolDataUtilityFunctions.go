package databaseAccess

// LoLPlayerCore
func (player *LoLPlayerCore) validate(leagueId, teamId, playerId int) (bool, string, error) {
	return validate(
		player.uniqueness(leagueId, teamId, playerId))
}

func (player *LoLPlayerCore) ValidateNew(leagueId, teamId int) (bool, string, error) {
	return player.validate(leagueId, teamId, 0)
}

func (player *LoLPlayerCore) ValidateEdit(leagueId, teamId, playerId int) (bool, string, error) {
	return player.validate(leagueId, teamId, playerId)
}

func (player *LoLPlayerCore) uniqueness(leagueId, teamId, playerId int) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		//TODO: implement this
		return true
	}
}
