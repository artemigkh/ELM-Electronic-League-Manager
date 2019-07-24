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
	NameTooLong               = "Name too long"
	NameTooShort              = "Name too short"
	TagTooLong                = "Tag too long"
	GameIdentifierTooLong     = "Game identifier too long"
	GameIdentifierTooShort    = "Game identifier too short"
	TagTooShort               = "Tag too short"
	PasswordTooLong           = "Password too long"
	PasswordTooShort          = "Password too short"
	DescriptionTooLong        = "Description too long"
	MarkdownTooLong           = "Markdown too long"
	LeagueNameInUse           = "League name already in use"
	EmailInUse                = "Email in use"
	InvalidGame               = "Game string not valid"
	LeaguePermissionsWrong    = "League cannot be invisible but have public join enabled"
	TimeOutOfOrder            = "Start time before end time"
	PeriodOutOfOrder          = "League starts before signup ends"
	AdminLackingPermissions   = "Administrator must have all permissions enabled"
	TeamsMustBeDifferent      = "The two teams in game must be different"
	TeamInGameDoesNotExist    = "A team in this game does not exist in this league"
	PlayerGameIdentifierInUse = "Player game Identifier already in use"
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
