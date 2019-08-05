package dataModel

// Each incoming data structure should implement this to check data correctness
type ValidateFunc func(*string, *error) bool

func validate(validators ...ValidateFunc) (bool, string, error) {
	valid := true
	problem := ""
	var err error
	for _, validator := range validators {
		valid = valid && validator(&problem, &err)
	}
	return valid, problem, err
}

const (
	MaxTagLength         = 5
	MaxNameLength        = 50
	MaxDescriptionLength = 500
	MaxMdLength          = 50000
	MaxPasswordLength    = 64
	MinInformationLength = 2
	MinPasswordLength    = 8
)

type DataProblem string

const (
	NameTooLong                       = "Name too long"
	NameTooShort                      = "Name too short"
	TagTooLong                        = "Tag too long"
	GameIdentifierTooLong             = "Game identifier too long"
	GameIdentifierTooShort            = "Game identifier too short"
	TagTooShort                       = "Tag too short"
	PasswordTooLong                   = "Password too long"
	PasswordTooShort                  = "Password too short"
	DescriptionTooLong                = "Description too long"
	MarkdownTooLong                   = "Markdown too long"
	LeagueNameInUse                   = "League name already in use"
	EmailInUse                        = "Email in use"
	EmailMalformed                    = "Email is malformed"
	InvalidGame                       = "Game string not valid"
	LeaguePermissionsWrong            = "League cannot be invisible but have public join enabled"
	TimeOutOfOrder                    = "Start time before end time"
	PeriodOutOfOrder                  = "League starts before signup ends"
	AdminLackingPermissions           = "Administrator must have all permissions enabled"
	TeamsMustBeDifferent              = "The two teams in game must be different"
	TeamInGameDoesNotExist            = "A team in this game does not exist in this league"
	PlayerGameIdentifierInUse         = "Player game Identifier already in use"
	ExternalIdentifierInUse           = "An entity in this league already has this external id"
	GameConflict                      = "A team in this game already has a game starting at this time"
	GameNotDuringLeague               = "This game start time is not during the league competition period"
	AvailabilityStartNotDuringLeague  = "This availability start time is not during the league competition period"
	AvailabilityEndNotDuringLeague    = "This availability end time is not during the league competition period"
	GameDoesNotContainTeams           = "The teams in this game report are not in this game"
	AvailabilityOutOfOrder            = "Availability start time must be before end time"
	InvalidWeekday                    = "Weekday must be one of 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'"
	InvalidTimezone                   = "Timezone offset cannot be more than 24 hours"
	InvalidHour                       = "Hour must be between 0 and 23 inclusive"
	InvalidMinute                     = "Minute must be between 0 and 59 inclusive"
	InvalidWeeklyAvailabilityDuration = "Weekly availability duration can not be longer than a week"
	TournamentTypeNotSupported        = "The specified tournament type is not supported"
	AvailabilitiesNotDuringLeague     = "New league times would result in an availability outside the league competition period"
	GamesNotDuringLeague              = "New league times would result in a game scheduled outside the league competition period"
)

var ValidGameStrings = [...]string{
	"genericsport",
	"basketball",
	"curling",
	"football",
	"hockey",
	"rugby",
	"soccer",
	"volleyball",
	"waterpolo",
	"genericesport",
	"csgo",
	"leagueoflegends",
	"overwatch",
}
