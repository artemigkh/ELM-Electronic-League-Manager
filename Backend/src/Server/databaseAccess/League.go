package databaseAccess

import "github.com/Masterminds/squirrel"

// League
func getLeagueSelector() squirrel.SelectBuilder {
	return psql.Select(
		"league_id",
		"name",
		"description",
		"game",
		"public_view",
		"public_join",
		"signup_start",
		"signup_end",
		"league_start",
		"league_end",
	).From("league")
}

// LeagueCore
func (league *LeagueCore) Validate(leagueId int) (bool, string, error) {
	return validate(
		league.name(),
		league.uniqueness(leagueId),
		league.description(),
		league.game(),
		league.permissions(),
		league.timestamps())
}

func (league *LeagueCore) name() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(league.Name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(league.Name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (league *LeagueCore) uniqueness(leagueId int) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, err := Leagues.IsNameInUse(leagueId, league.Name)
		*errorDest = err
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = LeagueNameInUse
		} else {
			valid = true
		}
		return valid
	}
}

func (league *LeagueCore) description() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if len(league.Description) > MaxDescriptionLength {
			*problemDest = DescriptionTooLong
			return false
		} else {
			return true
		}
	}
}

func (league *LeagueCore) game() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		for _, g := range ValidGameStrings {
			if g == league.Game {
				return true
			}
		}
		*problemDest = InvalidGame
		return false
	}
}

func (league *LeagueCore) permissions() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if league.PublicJoin && !league.PublicView {
			*problemDest = LeaguePermissionsWrong
			return false
		} else {
			return true
		}
	}
}

func (league *LeagueCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if league.SignupStart > league.SignupEnd ||
			league.LeagueStart > league.LeagueEnd {
			*problemDest = TimeOutOfOrder
		} else if league.LeagueStart < league.SignupEnd {
			*problemDest = PeriodOutOfOrder
		} else {
			valid = true
		}
		return valid
	}
}
