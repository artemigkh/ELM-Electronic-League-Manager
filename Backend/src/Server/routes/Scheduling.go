package routes

import (
	"Server/scheduler"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

type recurringAvailability struct {
	Weekday     string `json:"weekday"`
	Timezone    int    `json:"timezone"`
	Hour        int    `json:"hour"`
	Minute      int    `json:"minute"`
	Duration    int    `json:"duration"`
	Constrained bool   `json:"constrained"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

type recurringAvailabilityWithId struct {
	Id          int    `json:"id"`
	Weekday     string `json:"weekday"`
	Timezone    int    `json:"timezone"`
	Hour        int    `json:"hour"`
	Minute      int    `json:"minute"`
	Duration    int    `json:"duration"`
	Constrained bool   `json:"constrained"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

var weekdays = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}

var tournamentStringToEnum = map[string]int{
	"roundrobin":       scheduler.RoundRobin,
	"doubleroundrobin": scheduler.DoubleRoundRobin,
}

type scheduleRequest struct {
	TournamentType    string `json:"tournamentType"`
	RoundsPerWeek     int    `json:"roundsPerWeek"`
	ConcurrentGameNum int    `json:"concurrentGameNum"`
	GameDuration      int    `json:"gameDuration"`
}

type teamInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	IconSmall string `json:"iconSmall"`
}

type scheduledGame struct {
	Team1Id  int   `json:"team1Id"`
	Team2Id  int   `json:"team2Id"`
	GameTime int64 `json:"gameTime"`
}

type schedule struct {
	Teams []teamInfo      `json:"teams"`
	Games []scheduledGame `json:"games"`
}

/**
* @api{POST} /api/scheduling/recurringAvailability Add Recurring Scheduling Availability
* @apiName addAvailability
* @apiGroup Scheduling
* @apiDescription Provide a recurring time that is available for game scheduling in this league
*
* @apiParam {string} weekday 'sunday', 'monday', etc.
* @apiParam {number} timezone Timezone as offset in seconds east from UTC
* @apiParam {number} hour Number between 0-24 that is the start hour of the availability
* @apiParam {number} minute Number between 0-60 that is the start minute of the availability
* @apiParam {number} duration Duration of the availability in minutes
* @apiParam {boolean} constrained False if availability applies to every week of the league time
* @apiParam {number} start If constrained is True, the unix timestamp of the constraint start for this repeating availability
* @apiParam {number} end If allLeague is True, the unix timestamp of the constraint end for this repeating availability
*
* @apiSuccess {int} id The unique numerical identifier of the availability
*
* @apiError notLoggedIn No user is logged in
* @apiError noActiveLeague There is no active league selected
* @apiError notAdmin Currently logged in user is not a league administrator
 */
func addRecurringAvailability(ctx *gin.Context) {
	var availability recurringAvailability
	err := ctx.ShouldBindJSON(&availability)
	if checkJsonErr(ctx, err) {
		return
	}

	var day time.Weekday
	var ok bool
	if day, ok = weekdays[availability.Weekday]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalidDayOfWeek"})
		return
	}
	if failIfTimezoneOffsetTooLarge(ctx, math.Abs(float64(availability.Timezone))) {
		return
	}
	if failIfHourTooLarge(ctx, availability.Hour) {
		return
	}
	if failIfMinuteTooLarge(ctx, availability.Minute) {
		return
	}
	if failIfDurationTooLarge(ctx, availability.Duration) {
		return
	}

	id, err := LeaguesDAO.AddRecurringAvailability(ctx.GetInt("leagueId"), int(day),
		availability.Timezone, availability.Hour, availability.Minute, availability.Duration,
		availability.Constrained, availability.Start, availability.End)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

/**
* @api{PUT} /api/scheduling/recurringAvailability Edit Recurring Scheduling Availability
* @apiName editAvailability
* @apiGroup Scheduling
* @apiDescription Edit a recurring time that is available for game scheduling in this league
*
* @apiParam {number} id The unique numerical identifier of the availability
* @apiParam {string} weekday 'sunday', 'monday', etc.
* @apiParam {number} timezone Timezone as offset in seconds east from UTC
* @apiParam {number} hour Number between 0-24 that is the start hour of the availability
* @apiParam {number} minute Number between 0-60 that is the start minute of the availability
* @apiParam {number} duration Duration of the availability in minutes
* @apiParam {boolean} constrained False if availability applies to every week of the league time
* @apiParam {number} start If constrained is True, the unix timestamp of the constraint start for this repeating availability
* @apiParam {number} end If allLeague is True, the unix timestamp of the constraint end for this repeating availability
*
* @apiSuccess {int} id The unique numerical identifier of the availability
*
* @apiError notLoggedIn No user is logged in
* @apiError noActiveLeague There is no active league selected
* @apiError notAdmin Currently logged in user is not a league administrator
 */
func editRecurringAvailability(ctx *gin.Context) {
	var availability recurringAvailabilityWithId
	err := ctx.ShouldBindJSON(&availability)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfAvailabilityDoesNotExist(ctx, ctx.GetInt("leagueId"), availability.Id) {
		return
	}

	var day time.Weekday
	var ok bool
	if day, ok = weekdays[availability.Weekday]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalidDayOfWeek"})
		return
	}
	if failIfTimezoneOffsetTooLarge(ctx, math.Abs(float64(availability.Timezone))) {
		return
	}
	if failIfHourTooLarge(ctx, availability.Hour) {
		return
	}
	if failIfMinuteTooLarge(ctx, availability.Minute) {
		return
	}
	if failIfDurationTooLarge(ctx, availability.Duration) {
		return
	}

	err = LeaguesDAO.EditRecurringAvailability(ctx.GetInt("leagueId"), availability.Id, int(day),
		availability.Timezone, availability.Hour, availability.Minute, availability.Duration,
		availability.Constrained, availability.Start, availability.End)
	if checkErr(ctx, err) {
		return
	}
}

/**
* @api{DELETE} /api/scheduling/recurringAvailability/:id Remove Recurring Scheduling Availability
* @apiName removeAvailability
* @apiGroup Scheduling
* @apiDescription Remove a recurring time that is available for game scheduling in this league by id
*
* @apiParam {int} id the unique numerical identifier of the scheduling availability
*
* @apiError notLoggedIn No user is logged in
* @apiError noActiveLeague There is no active league selected
* @apiError notAdmin Currently logged in user is not a league administrator
 */
func deleteRecurringAvailability(ctx *gin.Context) {
	if failIfAvailabilityDoesNotExist(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId")) {
		return
	}

	err := LeaguesDAO.RemoveRecurringAvailabilities(ctx.GetInt("leagueId"), ctx.GetInt("urlId"))
	if checkErr(ctx, err) {
		return
	}
}

/**
 * @api{POST} /api/scheduling/schedule Get Schedule
 * @apiGroup Scheduling
 * @apiDescription Get a schedule of games for league given input parameters
 *

 * @apiParam {string} tournamentType 'roundrobin' or 'doubleroundrobin'
 * @apiParam {int} roundsPerWeek Where a round is where every team plays exactly once
 * @apiParam {int} concurrentGameNum How many games at once can be scheduled
 * @apiParam {int} gameDuration Duration of a game in minutes

 * @apiSuccess {jsonArray} teams An array of JSON objects, each representing a team
 * @apiSuccess {int} teams.id The unique numerical identifier of the team
 * @apiSuccess {string} teams.name The name of the team
 * @apiSuccess {string} teams.tag The tag of the team
 * @apiSuccess {string} teams.iconSmall The small icon filename
 * @apiSuccess {jsonArray} games An array of JSON objects, each representing a game
 * @apiSuccess {int} games.gameTime The unix time of when the game is scheduled for
 * @apiSuccess {int} games.team1Id The unique numerical identifier of team 1
 * @apiSuccess {int} games.team2Id The unique numerical identifier of team 2
 *
* @apiError notLoggedIn No user is logged in
* @apiError noActiveLeague There is no active league selected
* @apiError notAdmin Currently logged in user is not a league administrator
* @apiError schedulingProblem A schedule was not able to be generated given input parameters. Comes with
a description json of the form {error: "description"}
*/
func getSchedule(ctx *gin.Context) {
	var scheduleParams scheduleRequest
	err := ctx.ShouldBindJSON(&scheduleParams)
	if checkJsonErr(ctx, err) {
		return
	}

	var tournamentType int
	var ok bool
	if tournamentType, ok = tournamentStringToEnum[scheduleParams.TournamentType]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalidTournamentType"})
		return
	}
	if failIfDurationTooLarge(ctx, scheduleParams.GameDuration) {
		return
	}

	teamSummary, err := LeaguesDAO.GetTeamSummary(ctx.GetInt("leagueId"))
	var teamIds []int
	for _, team := range teamSummary {
		teamIds = append(teamIds, team.Id)
	}

	fmt.Printf("tournament type: %v \n", tournamentType)
	fmt.Printf("teamids: %v \n", teamIds)

	availabilities, err := LeaguesDAO.GetSchedulingAvailabilities(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	fmt.Printf("availabilities: %v \n", availabilities)

	leagueInformation, err := LeaguesDAO.GetLeagueInformation(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	fmt.Printf("start: %v, end: %v \n", leagueInformation.LeagueStart, leagueInformation.LeagueEnd)

	s := scheduler.Scheduler{}
	s.InitScheduler(tournamentType, scheduleParams.RoundsPerWeek, scheduleParams.ConcurrentGameNum,
		time.Duration(scheduleParams.GameDuration)*time.Minute,
		time.Unix(int64(leagueInformation.LeagueStart), 0), time.Unix(int64(leagueInformation.LeagueEnd), 0),
		teamIds)

	for _, avail := range availabilities {
		s.AddWeeklyAvailability(time.Weekday(avail.Weekday), avail.Hour, avail.Minute,
			time.Duration(avail.Duration)*time.Minute)
	}

	games, err := s.GetSchedule()
	if err != nil {
		print(err.Error())
	}

	var scheduledGames []scheduledGame
	for _, game := range games {
		fmt.Printf("%v vs %v - %v\n", game.Team1Id, game.Team2Id, time.Unix(game.GameTime, 0).Format(time.UnixDate))
		scheduledGames = append(scheduledGames, scheduledGame{
			Team1Id:  game.Team1Id,
			Team2Id:  game.Team2Id,
			GameTime: game.GameTime,
		})
	}

	var teamsInfo []teamInfo
	for _, team := range teamSummary {
		teamsInfo = append(teamsInfo, teamInfo{
			Id:        team.Id,
			Name:      team.Name,
			Tag:       team.Tag,
			IconSmall: team.IconSmall,
		})
	}

	ctx.JSON(http.StatusOK, schedule{
		Teams: teamsInfo,
		Games: scheduledGames,
	})
}

/**
 * @api{GET} /api/scheduling/availabilities/ Get League Scheduling Availabilities
 * @apiGroup Scheduling
 * @apiDescription Get time periods that games can be scheduled for this league
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing an availability
 * @apiSuccess {int} _.id The unique numerical identifier of the team
 * @apiSuccess {int} _.weekday The day of the week. 0 = sunday, 1 = monday, etc.
 * @apiSuccess {int} _.timezone The timezone of the availability in seconds east of UTC
 * @apiSuccess {int} _.hour The hour of of the availability
 * @apiSuccess {int} _.minute The minute of of the availability
 * @apiSuccess {int} _.duration The duration of of the availability
 * @apiSuccess {bool} _.constrained If the availability is constrained
 * @apiSuccess {int} _.start Unix time of constraint start
 * @apiSuccess {int} _.end Unix time of constraint end
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getAvailabilities(ctx *gin.Context) {
	availabilities, err := LeaguesDAO.GetSchedulingAvailabilities(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, availabilities)
}

func RegisterSchedulingHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())
	g.POST("/recurringAvailability", authenticate(), failIfNotLeagueAdmin(), addRecurringAvailability)
	g.PUT("/recurringAvailability", authenticate(), failIfNotLeagueAdmin(), editRecurringAvailability)
	g.DELETE("/recurringAvailability/:id", getUrlId(), authenticate(), failIfNotLeagueAdmin(), deleteRecurringAvailability)
	g.POST("/schedule", authenticate(), failIfNotLeagueAdmin(), getSchedule)
	g.GET("/availabilities", getAvailabilities)
}
