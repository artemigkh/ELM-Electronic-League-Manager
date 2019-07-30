package dataModel

type TeamDAO interface {
	// Teams
	CreateTeam(leagueId, userId int, teamInfo TeamCore) (int, error)
	CreateTeamWithIcon(leagueId, userId int, teamInfo TeamCore, iconSmall, iconLarge string) (int, error)
	CreateTeamWithPlayers(leagueId, userId int, teamInfo TeamCore, players []PlayerCore, iconSmall, iconLarge string) (int, error)
	DeleteTeam(teamId int) error
	UpdateTeam(teamId int, teamInformation TeamCore) error
	UpdateTeamIcon(teamId int, small, large string) error
	GetTeamInformation(teamId int) (*TeamWithPlayers, error)
	GetTeamWithRosters(teamId int) (*TeamWithRosters, error)
	GetAllTeamsInLeague(leagueId int) ([]*TeamWithPlayers, error)
	GetAllTeamsInLeagueWithRosters(leagueId int) ([]*TeamWithRosters, error)
	GetAllTeamDisplaysInLeague(leagueId int) ([]*TeamDisplay, error)

	// Players
	CreatePlayer(leagueId, teamId int, playerInfo PlayerCore) (int, error)
	DeletePlayer(playerId int) error
	UpdatePlayer(playerId int, playerInfo PlayerCore) error

	// Get Information For Team and Player Management
	GetTeamPermissions(teamId, userId int) (*TeamPermissionsCore, error)
	IsInfoInUse(leagueId, teamId int, name, tag string) (bool, string, error)
	DoesTeamExistInLeague(leagueId, teamId int) (bool, error)
	IsTeamActive(leagueId, teamId int) (bool, error)
	DoesPlayerExist(leagueId, teamId, playerId int) (bool, error)

	// Managers
	ChangeManagerPermissions(teamId, userId int, teamPermissionInformation TeamPermissionsCore) error
}

type TeamWithPlayersCore struct {
	Team    TeamCore     `json:"team"`
	Icon    string       `json:"icon"`
	Players []PlayerCore `json:"players"`
}

func (team *TeamWithPlayersCore) Validate(leagueId int, teamDao TeamDAO) (bool, string, error) {
	valid, problem, err := team.Team.ValidateNew(leagueId, teamDao)
	if !valid || problem != "" || err != nil {
		return valid, problem, err
	}

	// Check that each player is unique from the other non-existing players
	for i := 0; i < len(team.Players); i++ {
		otherPlayers := make([]*Player, 0)
		for j := 0; j < len(team.Players); j++ {
			if i != j {
				otherPlayers = append(otherPlayers, &Player{
					PlayerId:       0,
					Name:           team.Players[j].Name,
					GameIdentifier: team.Players[j].GameIdentifier,
					MainRoster:     team.Players[j].MainRoster,
				})
			}
		}

		valid, problem := team.Players[i].uniqueness(-1, otherPlayers)
		if !valid {
			return false, problem, nil
		}
	}

	// Validate each player normally
	for _, player := range team.Players {
		valid, problem, err := player.ValidateNew(leagueId, -1, teamDao)
		if !valid || problem != "" || err != nil {
			return valid, problem, err
		}
	}

	return true, "", nil
}

type TeamWithPlayers struct {
	TeamId      int       `json:"teamId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	IconSmall   string    `json:"iconSmall"`
	IconLarge   string    `json:"iconLarge"`
	Wins        int       `json:"wins"`
	Losses      int       `json:"losses"`
	Players     []*Player `json:"players"`
}

type TeamWithRosters struct {
	TeamId           int       `json:"teamId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Tag              string    `json:"tag"`
	IconSmall        string    `json:"iconSmall"`
	IconLarge        string    `json:"iconLarge"`
	Wins             int       `json:"wins"`
	Losses           int       `json:"losses"`
	MainRoster       []*Player `json:"mainRoster"`
	SubstituteRoster []*Player `json:"substituteRoster"`
}

type LoLTeamWithRosters struct {
	TeamId           int          `json:"teamId"`
	Name             string       `json:"name"`
	Description      string       `json:"description"`
	Tag              string       `json:"tag"`
	IconSmall        string       `json:"iconSmall"`
	IconLarge        string       `json:"iconLarge"`
	Wins             int          `json:"wins"`
	Losses           int          `json:"losses"`
	MainRoster       []*LoLPlayer `json:"mainRoster"`
	SubstituteRoster []*LoLPlayer `json:"substituteRoster"`
}

type TeamWithManagers struct {
	TeamId    int            `json:"teamId"`
	Name      string         `json:"name"`
	Tag       string         `json:"tag"`
	IconSmall string         `json:"iconSmall"`
	Managers  []*TeamManager `json:"managers"`
}

type TeamDisplay struct {
	TeamId    int    `json:"teamId"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	IconSmall string `json:"iconSmall"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
}

type TeamCore struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

func (team *TeamCore) validate(leagueId, teamId int, teamDao TeamDAO) (bool, string, error) {
	return validate(
		team.name(),
		team.uniqueness(leagueId, teamId, teamDao),
		team.tag(),
		team.description())
}

func (team *TeamCore) ValidateNew(leagueId int, teamDao TeamDAO) (bool, string, error) {
	return team.validate(leagueId, 0, teamDao)
}

func (team *TeamCore) ValidateEdit(leagueId, teamId int, teamDao TeamDAO) (bool, string, error) {
	return team.validate(leagueId, teamId, teamDao)
}

func (team *TeamCore) name() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(team.Name) > MaxNameLength {
			*problemDest = NameTooLong
		} else if len(team.Name) < MinInformationLength {
			*problemDest = NameTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (team *TeamCore) uniqueness(leagueId, teamId int, teamDao TeamDAO) ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, problem, err := teamDao.IsInfoInUse(leagueId, teamId, team.Name, team.Tag)
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = problem
		} else {
			valid = true
		}
		return valid
	}
}

func (team *TeamCore) tag() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(team.Tag) > MaxTagLength {
			*problemDest = TagTooLong
		} else if len(team.Tag) < MinInformationLength {
			*problemDest = TagTooShort
		} else {
			valid = true
		}
		return valid
	}
}

func (team *TeamCore) description() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if len(team.Description) > MaxDescriptionLength {
			*problemDest = DescriptionTooLong
			return false
		} else {
			return true
		}
	}
}
