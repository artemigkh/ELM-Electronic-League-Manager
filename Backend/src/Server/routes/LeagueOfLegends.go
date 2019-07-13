package routes

import (
	"Server/databaseAccess"
	"Server/lolApi"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/securecookie"
	"net/http"
)

//type TournamentCallback struct {
//	WinningTeam    []SummonerInformation `json:"winningTeam"`
//	LosingTeam     []SummonerInformation `json:"losingTeam"`
//	TournamentCode string                `json:"shortCode"`
//	GameId         string                `json:"gameId"`
//}
//
//// Team Stats Mappings
//const (
//	GameTime     = iota
//	FirstBloods  = iota
//	FirstTurrets = iota
//	TeamKDA      = iota
//)
//
//// Individual Stats Mappings
//const (
//	DamagePerMinute    = iota
//	GoldPerMinute      = iota
//	CsPerMinute        = iota
//	PlayerKDA          = iota
//	Kills              = iota
//	Deaths             = iota
//	Assists            = iota
//	VisionWardsPlaced  = iota
//	ControlWardsPlaced = iota
//)
//
//var tierOrder = map[string]int{
//	"IRON":        0,
//	"BRONZE":      1,
//	"SILVER":      2,
//	"GOLD":        3,
//	"PLATINUM":    4,
//	"DIAMOND":     5,
//	"MASTER":      6,
//	"GRANDMASTER": 7,
//	"CHALLENGER":  8,
//}

func leagueOfLegendsGetSummonerId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var playerInfo databaseAccess.LoLPlayerCore
		err := ctx.ShouldBindBodyWith(&playerInfo, binding.JSON)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
		}
		summonerId, err := LoLApi.GetSummonerId(playerInfo.GameIdentifier)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lolApiError"})
		}
		ctx.Set("externalId", summonerId)
		ctx.Next()
	}
}

func leagueOfLegendsGenerateExternalGameId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("externalGameId", hex.EncodeToString(securecookie.GenerateRandomKey(16)))
		ctx.Next()
	}
}

//TODO: validate created game struct
//TODO: add case for adding game result to series
func receiveCompletedTournamentGame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var matchInformation lolApi.MatchInformation
		if bindAndCheckErr(ctx, &matchInformation) {
			return
		}
	}

	//
	//return endpoint{
	//	Entity:     Player,
	//	AccessType: Edit,
	//	BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &matchInformation) },
	//	//IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
	//	//	return player.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), getPlayerId(ctx))
	//	//},
	//	Core: func(ctx *gin.Context) (interface{}, error) {
	//		fmt.Printf("%+v", matchInformation)
	//		return nil, nil
	//		//if err := GamesDAO.ReportGameByExternalId(matchInformation.GameId, databaseAccess.GameResult{
	//		//	WinnerId:   matchInformation.WinningTeamId,
	//		//	LoserId:    matchInformation.LosingTeamId,
	//		//	ScoreTeam1: func() int { if  matchInformation.WinningTeamId == matchInformation.Team1Id { return 1 } else { return 0 } }(),
	//		//	ScoreTeam2: func() int { if  matchInformation.WinningTeamId == matchInformation.Team2Id { return 1 } else { return 0 } }(),
	//		//}); err != nil {
	//		//	return nil, err
	//		//}
	//		//return nil, LeagueOfLegendsDAO.UpdateLoLPlayer(
	//		//	getPlayerId(ctx), getExternalId(ctx), player)
	//	},
	//}.createEndpointHandler()
}

//func tournamentCallback(ctx *gin.Context) {
//	var callbackInfo TournamentCallback
//	err := ctx.ShouldBindJSON(&callbackInfo)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	// Get team information from database
//	gameInfo, err := GamesDAO.GetGameInformationFromExternalId(callbackInfo.TournamentCode)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	team1Info, err := TeamsDAO.GetTeamInformation(gameInfo.Team1Id)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	// Find which team id won and which team id lost by querying a member
//	var winningId int
//	var losingId int
//	var team1Score int
//	var team2Score int
//
//	losingId = gameInfo.Team1Id
//	for _, player := range team1Info.Players {
//		if callbackInfo.WinningTeam[0].SummonerId == player.ExternalId {
//			winningId = gameInfo.Team1Id
//			losingId = gameInfo.Team2Id
//
//			team1Score = 1
//			team2Score = 0
//		}
//	}
//	if losingId == gameInfo.Team1Id {
//		winningId = gameInfo.Team2Id
//
//		team1Score = 0
//		team2Score = 1
//	}
//
//	// Report Game Complete
//	if failIfGameDoesNotExist(ctx, gameInfo.Id) {
//		return
//	}
//	if failIfGameDoesNotContainWinner(ctx, gameInfo.Id, winningId) {
//		return
//	}
//
//	//report the result
//	err = GamesDAO.ReportGame(
//		databaseAccess.GameDTO{
//			Id:         gameInfo.Id,
//			LeagueId:   gameInfo.LeagueId,
//			Complete:   false,
//			WinnerId:   winningId,
//			LoserId:    losingId,
//			ScoreTeam1: team1Score,
//			ScoreTeam2: team2Score,
//		})
//	if checkErr(ctx, err) {
//		return
//	}
//
//	// Get match stats from LoL Api
//	matchStats, err := LoLApi.GetMatchStats(callbackInfo.GameId)
//	if err != nil {
//		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lolApiError"})
//	}
//
//	err = LeagueOfLegendsDAO.ReportEndGameStats(gameInfo.LeagueId, gameInfo.Id, winningId, losingId, matchStats)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	print(fmt.Sprintf("%+v", matchStats))
//}
//
//func getPlayerStats(ctx *gin.Context) {
//	playerStats, err := LeagueOfLegendsDAO.GetPlayerStats(ctx.GetInt("getLeagueId"))
//	if checkErr(ctx, err) {
//		return
//	}
//	ctx.JSON(http.StatusOK, playerStats)
//}
//
//func getTeamStats(ctx *gin.Context) {
//	teamStats, err := LeagueOfLegendsDAO.GetTeamStats(ctx.GetInt("getLeagueId"))
//	if checkErr(ctx, err) {
//		return
//	}
//	ctx.JSON(http.StatusOK, teamStats)
//}
//
//func getChampionStats(ctx *gin.Context) {
//	championStats, err := LeagueOfLegendsDAO.GetChampionStats(ctx.GetInt("getLeagueId"))
//	if checkErr(ctx, err) {
//		return
//	}
//	ctx.JSON(http.StatusOK, championStats)
//}
//

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createLoLPlayer
func createNewLoLPlayer() gin.HandlerFunc {
	var player databaseAccess.LoLPlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Create,
		BindData:   func(ctx *gin.Context) bool { return bindRepeatedAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateNew(getLeagueId(ctx), getTeamId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			playerId, err := LeagueOfLegendsDAO.CreateLoLPlayer(
				getLeagueId(ctx), getTeamId(ctx), getExternalId(ctx), player)
			return gin.H{"playerId": playerId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateLoLPlayer
func updateLoLPlayer() gin.HandlerFunc {
	var player databaseAccess.LoLPlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Edit,
		BindData:   func(ctx *gin.Context) bool { return bindRepeatedAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), getPlayerId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, LeagueOfLegendsDAO.UpdateLoLPlayer(
				getPlayerId(ctx), getExternalId(ctx), player)
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLoLTeamWithRosters
func getLoLTeamWithRosters() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			teamStub, err := LeagueOfLegendsDAO.GetLoLTeamStub(getTeamId(ctx))
			if checkErr(ctx, err) {
				return nil, err
			}
			return LoLApi.CompletePlayerStubs(teamStub)
		},
	}.createEndpointHandler()
}

func RegisterLeagueOfLegendsHandlers(g *gin.RouterGroup) {
	g.POST("/receiveCompletedTournamentGame", receiveCompletedTournamentGame())
	//g.Use(getActiveLeague())
	//
	//g.GET("/stats/player", getPlayerStats)
	//g.GET("/stats/team", getTeamStats)
	//g.GET("/stats/champion", getChampionStats)
	g.POST("/games", leagueOfLegendsGenerateExternalGameId(), createNewGame())
	withTeamId := g.Group("/teams/:teamId", storeTeamId())
	withTeamId.POST("/players", leagueOfLegendsGetSummonerId(), createNewLoLPlayer())
	withTeamId.PUT("/players/:playerId", storePlayerId(), leagueOfLegendsGetSummonerId(), updateLoLPlayer())
	withTeamId.GET("/withRosters", getLoLTeamWithRosters())
	//g.PUT("/teams/updatePlayer", authenticate(), leagueOfLegendsGetSummonerId(), updatePlayer)
	//g.GET("/teams/:id", storeUrlId(), leagueOfLegendsGetTeamInformation)
	//g.GET("/tournamentCode/:id", storeUrlId(), getTournamentCodeForGame)

}
