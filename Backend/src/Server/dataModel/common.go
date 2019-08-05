package dataModel

import "github.com/badoux/checkmail"

func validateName(name string) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func validateGameIdentifier(gameIdentifier string) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(gameIdentifier) > MaxNameLength {
			*problemDest = GameIdentifierTooLong
		} else if len(gameIdentifier) < MinInformationLength {
			*problemDest = GameIdentifierTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func validateAvailabilityTimestamps(start, end int) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if start > end {
			*problemDest = AvailabilityOutOfOrder
			return false
		} else {
			return true
		}
	}
}

func validateDuringLeague(leagueId, epochTime int, leagueDao LeagueDAO, onFailMessage string) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		leagueInformation, err := leagueDao.GetLeagueInformation(leagueId)
		if err != nil {
			*errorDest = err
		} else if epochTime < leagueInformation.LeagueStart ||
			epochTime > leagueInformation.LeagueEnd {
			*problemDest = onFailMessage
		} else {
			valid = true
		}
		return valid
	}
}

func validateEmailFormat(email string) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if err := checkmail.ValidateFormat(email); err != nil {
			*problemDest = EmailMalformed
			return false
		} else {
			return true
		}
	}
}

func validatePasswordLength(password string) ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(password) > MaxPasswordLength {
			*problemDest = PasswordTooLong
		} else if len(password) < MinPasswordLength {
			*problemDest = PasswordTooShort
		} else {
			valid = true
		}
		return valid
	}
}
