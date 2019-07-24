package dataModel

type AvailabilityCore struct {
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`
}

type Availability struct {
	AvailabilityId int `json:"availabilityId"`
	StartTime      int `json:"startTime"`
	EndTime        int `json:"endTime"`
}

func (avail *AvailabilityCore) Validate(leagueId int) (bool, string, error) {
	//TODO: implement these validate functions
	return validate(avail.timestamps())
}

func (avail *AvailabilityCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
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

func (avail *WeeklyAvailabilityCore) validate(leagueId, availabilityId int) (bool, string, error) {
	//TODO: implement these validate functions
	return validate(avail.timestamps())
}

func (avail *WeeklyAvailabilityCore) ValidateNew(leagueId int) (bool, string, error) {
	return avail.validate(leagueId, 0)
}

func (avail *WeeklyAvailabilityCore) ValidateEdit(leagueId, availabilityId int) (bool, string, error) {
	return avail.validate(leagueId, availabilityId)
}

func (avail *WeeklyAvailabilityCore) timestamps() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
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

func (params *SchedulingParameters) Validate() (bool, string, error) {
	//TODO: implement these validate functions
	return validate(params.tournamentType())
}

func (params *SchedulingParameters) tournamentType() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		return true
	}
}
