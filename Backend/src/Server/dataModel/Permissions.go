package dataModel

type LeaguePermissionsCore struct {
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

func (p *LeaguePermissionsCore) Validate() (bool, string, error) {
	return validate(p.consistent())
}

func (p *LeaguePermissionsCore) consistent() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if (p.CreateTeams || p.EditGames || p.EditTeams) && p.Administrator {
			*problemDest = AdminLackingPermissions
			return false
		} else {
			return true
		}
	}
}

type TeamPermissionsCore struct {
	Administrator bool `json:"administrator"`
	Information   bool `json:"information"`
	Games         bool `json:"games"`
}

func (p *TeamPermissionsCore) Validate() (bool, string, error) {
	return validate(p.consistent())
}

func (p *TeamPermissionsCore) consistent() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if (p.Information || p.Games) && p.Administrator {
			*problemDest = AdminLackingPermissions
			return false
		} else {
			return true
		}
	}
}

type TeamPermissions struct {
	TeamId        int    `json:"teamId"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	IconSmall     string `json:"iconSmall"`
	Administrator bool   `json:"administrator"`
	Information   bool   `json:"information"`
	Games         bool   `json:"games"`
}
