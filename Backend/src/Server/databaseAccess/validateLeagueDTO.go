package databaseAccess

const (
	MaxDescriptionLength = 500
)

func leagueNameUniquenessValid(leagueId int, name string, problem *string, errorDest *error) bool {
	inUse, err := leaguesDAO.IsNameInUse(leagueId, name)
	*errorDest = err
	if err != nil {
		*errorDest = err
		return false
	} else if inUse {
		*problem = "league name already in use"
		return false
	} else {
		return true
	}
}

func descriptionStringValid(desc string, problem *string) bool {
	valid := false
	if len(desc) > MaxDescriptionLength {
		*problem = "description too long"
	} else {
		valid = true
	}
	return valid
}

func basicPermissionsValid(publicView, publicJoin bool, problem *string) bool {
	valid := false
	if publicJoin && !publicView {
		*problem = "league cannot be invisible but have public join enabled"
	} else {
		valid = true
	}
	return valid
}

func timestampsValid(signupStart, signupEnd, leagueStart, leagueEnd int, problem *string) bool {
	valid := false
	if signupStart < signupEnd || leagueStart < leagueEnd {
		*problem = "period start times before end times"
	} else if leagueStart < signupEnd {
		*problem = "league starts before signup ends"
	} else {
		valid = true
	}
	return valid
}

func (d *DTOValidator) ValidateLeagueDTO(league LeagueDTO) (bool, string, error) {
	var err error
	problem := ""
	valid := nameStringValid(league.Name, &problem) &&
		leagueNameUniquenessValid(league.Id, league.Name, &problem, &err) &&
		descriptionStringValid(league.Description, &problem) &&
		basicPermissionsValid(league.PublicView, league.PublicJoin, &problem) &&
		timestampsValid(league.SignupStart,
			league.SignupEnd,
			league.LeagueStart,
			league.LeagueEnd,
			&problem)

	return valid, problem, err
}
