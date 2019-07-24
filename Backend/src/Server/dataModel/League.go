package dataModel

type LeagueDAO interface {
	// Modify League
	CreateLeague(userId int, leagueInfo LeagueCore) (int, error)
	UpdateLeague(leagueId int, leagueInfo LeagueCore) error
	JoinLeague(leagueId, userId int) error

	// Permissions
	SetLeaguePermissions(leagueId, userId int, permissions LeaguePermissionsCore) error
	//GetLeaguePermissions(leagueId, userId int) (*LeaguePermissionsDTO, error)
	GetTeamManagerInformation(leagueId int) ([]*TeamWithManagers, error)
	IsLeagueViewable(leagueId, userId int) (bool, error)
	CanJoinLeague(leagueId, userId int) (bool, error)

	// Get Information About Leagues
	DoesLeagueExist(leagueId int) (bool, error)
	GetLeagueInformation(leagueId int) (*League, error)
	IsNameInUse(leagueId int, name string) (bool, error)
	GetPublicLeagueList() ([]*League, error)

	// Markdown
	GetMarkdownFile(leagueId int) (string, error)
	SetMarkdownFile(leagueId int, fileName string) error

	// Availabilities
	AddAvailability(leagueId int, availability AvailabilityCore) (int, error)
	GetAvailabilities(leagueId int) ([]*Availability, error)
	DeleteAvailability(availabilityId int) error

	AddWeeklyAvailability(leagueId int, availability WeeklyAvailabilityCore) (int, error)
	GetWeeklyAvailabilities(leagueId int) ([]*WeeklyAvailability, error)
	EditWeeklyAvailability(availabilityId int, availability WeeklyAvailabilityCore) error
	DeleteWeeklyAvailability(availabilityId int) error

	DoesAvailabilityExistInLeague(leagueId, availabilityId int) (bool, error)
}

type LeagueCore struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Game        string `json:"game"`
	PublicView  bool   `json:"publicView"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signupStart"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
}

type League struct {
	LeagueId    int    `json:"leagueId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Game        string `json:"game"`
	PublicView  bool   `json:"publicView"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signupStart"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
}

func (league *LeagueCore) validate(leagueId int, leagueDao LeagueDAO) (bool, string, error) {
	return validate(
		league.name(),
		league.uniqueness(leagueId, leagueDao),
		league.description(),
		league.game(),
		league.permissions(),
		league.timestamps())
}

func (league *LeagueCore) ValidateNew(leagueDao LeagueDAO) (bool, string, error) {
	return league.validate(0, leagueDao)
}

func (league *LeagueCore) ValidateEdit(leagueId int, leagueDao LeagueDAO) (bool, string, error) {
	return league.validate(leagueId, leagueDao)
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

func (league *LeagueCore) uniqueness(leagueId int, leagueDao LeagueDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, err := leagueDao.IsNameInUse(leagueId, league.Name)
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

type Markdown struct {
	Markdown string `json:"markdown"`
}

func (md *Markdown) Validate() (bool, string, error) {
	return validate(md.markdown())
}

func (md *Markdown) markdown() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if len(md.Markdown) > MaxMdLength {
			*problemDest = MarkdownTooLong
			return false
		} else {
			return true
		}
	}
}
