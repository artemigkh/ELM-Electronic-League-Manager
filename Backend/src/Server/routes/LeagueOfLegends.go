package routes

import (
	"Server/config"
	lolApi "github.com/artemigkh/GoLang-LeagueOfLegendsAPIV4Framework"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

/**
 * @api{POST} /api/league-of-legends/teams/addPlayer Add Player To Team
 * @apiGroup Teams
 * @apiDescription Create a new player and add him to the teams roster
 *
 * @apiParam {int} teamId The unique numerical identifier of the team the player is to be added to
 * @apiParam {string} name The name of the player (can be left blank)
 * @apiParam {string} gameIdentifier The in-game name identifier of the player (jersey number, ign, etc.)
 * @apiParam {bool} mainRoster If true, this player is on the main roster, otherwise is a substitute
 *
 * @apiSuccess {int} id the unique numerical identifier of the created player
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
 * @apiError gameIdentifierTooLong The game identifier exceeds 50 characters
 * @apiError nameTooLong The name exceeds 50 characters
 * @apiError gameIdentifierTooShort The game identifier is smaller than 2 characters
 * @apiError gameIdentifierInUse This game identifier is already in use in this league
 * @apiError lolApiError There was an error retrieving information from the league of legends api
 */
func leagueOfLegendsGetSummonerId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var playerInfo PlayerInformation
		err := ctx.ShouldBindBodyWith(&playerInfo, binding.JSON)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
		}
		summonerId, err := lolApi.FromName(playerInfo.Name).SummonerId()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lolApiError"})
		}
		ctx.Set("externalId", summonerId)
		ctx.Next()
	}
}

func RegisterLeagueOfLegendsHandlers(g *gin.RouterGroup, conf config.Config) {
	lolApi.Init("NA1", conf.GetLeagueOfLegendsApiKey())
	g.Use(getActiveLeague())

	g.POST("/teams/addPlayer", authenticate(), leagueOfLegendsGetSummonerId(), addPlayerToTeam)
	g.PUT("/teams/updatePlayer", authenticate(), leagueOfLegendsGetSummonerId(), updatePlayer)
}
