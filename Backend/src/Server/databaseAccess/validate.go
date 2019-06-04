package databaseAccess

// Each incoming data structure should implement this to check data correctness
type ValidateFunc func(*string, *error) bool
type Validator interface {
	Validate(int) (bool, string, error)
}

// checks data to be within valid bounds and logically consistent
// wraps each structs validate() method so that it can be mocked
// during unit testing
type StructValidatorWrapper interface {
	ValidateData(s Validator) (bool, string, error)
}
type StructValidator struct{}

func (v *StructValidator) ValidateNew(s Validator) (bool, string, error) {
	return s.Validate(0)
}
func (v *StructValidator) ValidateEdit(s Validator, id int) (bool, string, error) {
	return s.Validate(id)
}

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
	MaxNameLength        = 50
	MaxDescriptionLength = 500
	MinInformationLength = 2
)

type DataProblem string

const (
	NameTooLong            = "Name too long"
	NameTooShort           = "Name too short"
	DescriptionTooLong     = "Description too long"
	LeagueNameInUse        = "League name already in use"
	InvalidGame            = "Game string not valid"
	LeaguePermissionsWrong = "League cannot be invisible but have public join enabled"
	TimeOutOfOrder         = "Start time before end time"
	PeriodOutOfOrder       = "League starts before signup ends"
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

// LeaguePermissionsCore

// TeamPermissionsCore

// UserCreationInformation

// Markdown

// TeamCore

// PlayerCore

// GameCreationInformation

// GameTime

// GameResult

// AvailabilityCore

// WeeklyAvailabilityCore

// SchedulingParameters
