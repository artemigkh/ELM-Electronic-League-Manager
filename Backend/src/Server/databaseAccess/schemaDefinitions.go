package databaseAccess

// Permissions
type LeaguePermissionsCore struct {
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

type TeamPermissionsCore struct {
	Administrator bool `json:"administrator"`
	Information   bool `json:"information"`
	Games         bool `json:"games"`
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

// Users
type UserCreationInformation struct {
	Email    string
	Password string
}

type User struct {
	UserId            int                   `json:"userId"`
	Email             string                `json:"email"`
	LeaguePermissions *LeaguePermissionsCore `json:"leaguePermissions"`
	TeamPermissions   []*TeamPermissions    `json:"teamPermissions"`
}

type TeamManager struct {
	UserId        int    `json:"userId"`
	Email         string `json:"email"`
	Administrator bool   `json:"administrator"`
	Information   bool   `json:"information"`
	Games         bool   `json:"games"`
}

// Leagues
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

// Teams
type TeamCore struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
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
}

// Players
type PlayerCore struct {
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}

type Player struct {
	PlayerId       int    `json:"playerId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"`
	MainRoster     bool   `json:"mainRoster"`
}

// Games
type GameTime struct {
	GameTime int `json:"gameTime"`
}

type GameCreationInformation struct {
	Team1Id  int `json:"team1Id"`
	Team2Id  int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

type GameCore struct {
	GameTime int         `json:"gameTime"`
	Team1    TeamDisplay `json:"team1"`
	Team2    TeamDisplay `json:"team2"`
}

type GameResult struct {
	WinnerId   int `json:"winnerId"`
	LoserId    int `json:"loserId"`
	ScoreTeam1 int `json:"scoreTeam1"`
	ScoreTeam2 int `json:"scoreTeam2"`
}

type Game struct {
	GameId     int         `json:"gameId"`
	GameTime   int         `json:"gameTime"`
	Team1      TeamDisplay `json:"team1"`
	Team2      TeamDisplay `json:"team2"`
	WinnerId   int         `json:"winnerId"`
	LoserId    int         `json:"loserId"`
	ScoreTeam1 int         `json:"scoreTeam1"`
	ScoreTeam2 int         `json:"scoreTeam2"`
	Complete   bool        `json:"complete"`
}

// Scheduling
type AvailabilityCore struct {
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`
}

type Availability struct {
	AvailabilityId int `json:"availabilityId"`
	StartTime      int `json:"startTime"`
	EndTime        int `json:"endTime"`
}

type WeeklyAvailabilityCore struct {
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
	Weekday   string `json:"weekday"`
	Timezone  int    `json:"timezone"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	Duration  int    `json:"duration"`
}

type WeeklyAvailability struct {
	AvailabilityId int    `json:"availabilityId"`
	StartTime      int    `json:"startTime"`
	EndTime        int    `json:"endTime"`
	Weekday        string `json:"weekday"`
	Timezone       int    `json:"timezone"`
	Hour           int    `json:"hour"`
	Minute         int    `json:"minute"`
	Duration       int    `json:"duration"`
}

type SchedulingParameters struct {
	TournamentType    string `json:"tournamentType"`
	RoundsPerWeek     int    `json:"roundsPerWeek"`
	ConcurrentGameNum int    `json:"concurrentGameNum"`
	GameDuration      int    `json:"gameDuration"`
}

// Misc
type Markdown struct {
	Markdown string `json:"markdown"`
}
