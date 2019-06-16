package routes

import (
	"Server/databaseAccess"
	"Server/scheduler"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createAvailability
func createNewAvailability() gin.HandlerFunc {
	var availability databaseAccess.AvailabilityCore
	return endpoint{
		Entity:        Availability,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &availability) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return availability.Validate(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			availabilityId, err := LeaguesDAO.AddAvailability(getLeagueId(ctx), availability)
			return gin.H{"availabilityId": availabilityId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getAvailabilities
func getAvailabilities() gin.HandlerFunc {
	return endpoint{
		Entity:     Availability,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return LeaguesDAO.GetAvailabilities(getLeagueId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deleteAvailability
func deleteAvailability() gin.HandlerFunc {
	return endpoint{
		Entity:     Availability,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, LeaguesDAO.DeleteAvailability(getAvailabilityId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createWeeklyAvailability
func createNewWeeklyAvailability() gin.HandlerFunc {
	var availability databaseAccess.WeeklyAvailabilityCore
	return endpoint{
		Entity:        Availability,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &availability) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return availability.ValidateNew(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			availabilityId, err := LeaguesDAO.AddWeeklyAvailability(getLeagueId(ctx), availability)
			return gin.H{"availabilityId": availabilityId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getWeeklyAvailabilities
func getWeeklyAvailabilities() gin.HandlerFunc {
	return endpoint{
		Entity:     Availability,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return LeaguesDAO.GetWeeklyAvailabilities(getLeagueId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deleteWeeklyAvailability
func deleteWeeklyAvailability() gin.HandlerFunc {
	return endpoint{
		Entity:     Availability,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, LeaguesDAO.DeleteWeeklyAvailability(getAvailabilityId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/editWeeklyAvailability
func editWeeklyAvailability() gin.HandlerFunc {
	var availability databaseAccess.WeeklyAvailabilityCore
	return endpoint{
		Entity:     Availability,
		AccessType: Create,
		BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &availability) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return availability.ValidateEdit(getLeagueId(ctx), getAvailabilityId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, LeaguesDAO.EditWeeklyAvailability(getAvailabilityId(ctx), availability)
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/generateSchedule
func generateSchedule() gin.HandlerFunc {
	var schedulingParameters databaseAccess.SchedulingParameters
	return endpoint{
		Entity:        Availability,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &schedulingParameters) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return schedulingParameters.Validate() },
		Core: func(ctx *gin.Context) (interface{}, error) {
			// Get list of all teams and create map of team id to its display information
			var teamIds []int
			teamDisplay := make(map[int]databaseAccess.TeamDisplay)

			teams, err := TeamsDAO.GetAllTeamDisplaysInLeague(getLeagueId(ctx))
			if err != nil {
				return nil, err
			}

			for _, team := range teams {
				teamIds = append(teamIds, team.TeamId)
				teamDisplay[team.TeamId] = *team
			}

			fmt.Printf("teamids: %v \n", teamIds)

			// Get Availabilities for use in scheduler
			// TODO: add non-weekly availabilities as well here
			availabilities, err := LeaguesDAO.GetWeeklyAvailabilities(getLeagueId(ctx))
			if err != nil {
				return nil, err
			}
			fmt.Printf("availabilities: %v \n", availabilities)

			// Get competition start and end times for use in scheduler
			leagueInformation, err := LeaguesDAO.GetLeagueInformation(getLeagueId(ctx))
			if err != nil {
				return nil, err
			}
			fmt.Printf("start: %v, end: %v \n", leagueInformation.LeagueStart, leagueInformation.LeagueEnd)

			// Use Scheduler to generate list of games
			s := scheduler.Scheduler{}
			s.InitScheduler(
				s.GetTournamentFromString(schedulingParameters.TournamentType),
				schedulingParameters.RoundsPerWeek,
				schedulingParameters.ConcurrentGameNum,
				time.Duration(schedulingParameters.GameDuration)*time.Minute,
				time.Unix(int64(leagueInformation.LeagueStart), 0),
				time.Unix(int64(leagueInformation.LeagueEnd), 0),
				teamIds)

			for _, availability := range availabilities {
				s.AddWeeklyAvailability(
					s.GetWeekdayFromString(availability.Weekday),
					availability.Hour,
					availability.Minute,
					time.Duration(availability.Duration)*time.Minute)
			}

			scheduledGames, err := s.GetSchedule()
			if err != nil {
				print(err.Error())
				return nil, err
			}

			fmt.Printf("games: %+v\n", scheduledGames)

			var games []databaseAccess.GameCore
			for _, game := range scheduledGames {
				games = append(games, databaseAccess.GameCore{
					GameTime: game.GameTime,
					Team1:    teamDisplay[game.Team1Id],
					Team2:    teamDisplay[game.Team2Id],
				})
			}

			return games, nil
		},
	}.createEndpointHandler()
}

func RegisterSchedulingHandlers(g *gin.RouterGroup) {
	availabilities := g.Group("/availabilities")
	availabilities.POST("", createNewAvailability())
	availabilities.GET("", getAvailabilities())
	availabilities.DELETE("/:availabilityId", storeAvailabilityId(), deleteAvailability())

	weeklyAvailabilities := g.Group("/weeklyAvailabilities")
	weeklyAvailabilities.POST("", createNewWeeklyAvailability())
	weeklyAvailabilities.GET("", getWeeklyAvailabilities())
	weeklyAvailabilities.DELETE("/:availabilityId", storeAvailabilityId(), deleteWeeklyAvailability())
	weeklyAvailabilities.PUT("/:availabilityId", storeAvailabilityId(), editWeeklyAvailability())

	g.POST("/schedule", generateSchedule())
}
