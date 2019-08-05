package routes

import (
	"Server/dataModel"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"net/http"
	"strconv"
)

func leagueOfLegendsGetSummonerId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var playerInfo dataModel.LoLPlayerCore
		err := ctx.ShouldBindBodyWith(&playerInfo, binding.JSON)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorDescription": "malformedInput"})
			return
		}
		summonerId, err := LoLApi.GetSummonerId(playerInfo.GameIdentifier)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorDescription": fmt.Sprintf("Could not find information for summoner \"%v\"", playerInfo.GameIdentifier)})
			return
		}
		ctx.Set("externalId", summonerId)
		ctx.Next()
	}
}

//TODO: add case for adding game result to series
func receiveCompletedTournamentGame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var matchInformation dataModel.LoLMatchInformation
		if bindAndCheckErr(ctx, &matchInformation) {
			return
		}

		gameResult := dataModel.GameResult{
			WinnerId: matchInformation.WinningTeamId,
			LoserId:  matchInformation.LosingTeamId,
			ScoreTeam1: func() int {
				if matchInformation.WinningTeamId == matchInformation.Team1Id {
					return 1
				} else {
					return 0
				}
			}(),
			ScoreTeam2: func() int {
				if matchInformation.WinningTeamId == matchInformation.Team2Id {
					return 1
				} else {
					return 0
				}
			}(),
		}
		valid, problem, err := gameResult.ValidateByExternalId(matchInformation.GameId, GameDAO)
		if DataInvalid(ctx, valid, problem, err) {
			return
		}
		leagueId, gameId, err := GameDAO.ReportGameByExternalId(matchInformation.GameId, gameResult)
		if checkErr(ctx, err) {
			return
		}
		if err := LeagueOfLegendsDAO.ReportEndGameStats(
			leagueId, gameId, &matchInformation); checkErr(ctx, err) {
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func getPlayerStats(ctx *gin.Context) {
	playerStats, err := LeagueOfLegendsDAO.GetPlayerStats(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, playerStats)
}

func getTeamStats(ctx *gin.Context) {
	teamStats, err := LeagueOfLegendsDAO.GetTeamStats(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, teamStats)
}

func getChampionStats(ctx *gin.Context) {
	championStats, err := LeagueOfLegendsDAO.GetChampionStats(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, championStats)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createLoLPlayer
func createNewLoLPlayer() gin.HandlerFunc {
	var player dataModel.LoLPlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Create,
		BindData:   func(ctx *gin.Context) bool { return bindRepeatedAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateNew(getLeagueId(ctx), getTeamId(ctx), LeagueOfLegendsDAO)
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
	var player dataModel.LoLPlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Edit,
		BindData:   func(ctx *gin.Context) bool { return bindRepeatedAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), getPlayerId(ctx), LeagueOfLegendsDAO)
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
			team, err := LeagueOfLegendsDAO.GetLoLTeamStub(getTeamId(ctx))
			if checkErr(ctx, err) {
				return nil, err
			}
			return LoLApi.CompletePlayerStubs(team)
		},
	}.createEndpointHandler()
}

func createNewLoLTeamWithPlayers() gin.HandlerFunc {
	var team dataModel.LoLTeamWithPlayersCore
	return endpoint{
		Entity:     Team,
		AccessType: Create,
		BindData: func(ctx *gin.Context) bool {
			team = dataModel.LoLTeamWithPlayersCore{}
			return bindAndCheckErr(ctx, &team)
		},
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return team.Validate(getLeagueId(ctx), TeamDAO)
		},
		CustomResCore: func(ctx *gin.Context) {
			var err error

			for _, player := range team.Players {
				externalId, err := LoLApi.GetSummonerId(player.GameIdentifier)
				if checkErr(ctx, err) {
					return
				}
				player.ExternalId = externalId

				valid, problem, err := player.ValidateNew(getLeagueId(ctx), 0, LeagueOfLegendsDAO)
				if DataInvalid(ctx, valid, problem, err) {
					return
				}
			}

			smallIcon := ""
			largeIcon := ""
			if len(team.Icon) > 0 {
				smallIcon, largeIcon, err = IconManager.StoreNewIconFromBase64String(team.Icon)
				if checkErr(ctx, err) {
					return
				}
			}

			teamId, err := LeagueOfLegendsDAO.CreateLoLTeamWithPlayers(
				getLeagueId(ctx), getUserId(ctx), team.Team, team.Players, smallIcon, largeIcon)
			ctx.JSON(http.StatusCreated, gin.H{"teamId": teamId})
		},
	}.createEndpointHandler()
}

func registerTournament() gin.HandlerFunc {
	return endpoint{
		Entity:     League,
		AccessType: Edit,
		CustomResCore: func(ctx *gin.Context) {
			leagueInfo, err := LeagueDAO.GetLeagueInformation(getLeagueId(ctx))
			if checkErr(ctx, err) {
				return
			}

			hasRegisteredTournament, err := LeagueOfLegendsDAO.LeagueHasRegisteredTournament(getLeagueId(ctx))
			if checkErr(ctx, err) {
				return
			} else if hasRegisteredTournament {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": "Tournament already registered"})
				return
			}

			providerId, tournamentId, err := LoLTournamentApi.RegisterTournament(getLeagueId(ctx), "NA", leagueInfo.Name)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": "Tournament already registered"})
				return
			}

			if err = LeagueOfLegendsDAO.RegisterTournamentProvider(getLeagueId(ctx), providerId, tournamentId); checkErr(ctx, err) {
				return
			}

			ctx.Status(http.StatusOK)
		},
	}.createEndpointHandler()
}

func getTournamentCode() gin.HandlerFunc {
	return endpoint{
		Entity:     Game,
		AccessType: Edit,
		CustomResCore: func(ctx *gin.Context) {
			newCodeParam := ctx.DefaultQuery("new", "false")
			newCode, err := strconv.ParseBool(newCodeParam)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": "newParamMustBeBoolean"})
				return
			}

			gameInfo, err := GameDAO.GetGameInformation(getGameId(ctx))
			if checkErr(ctx, err) {
				return
			}

			team1Info, err := LeagueOfLegendsDAO.GetLoLTeamStub(gameInfo.Team1.TeamId)
			if checkErr(ctx, err) {
				return
			}

			team2Info, err := LeagueOfLegendsDAO.GetLoLTeamStub(gameInfo.Team2.TeamId)
			if checkErr(ctx, err) {
				return
			}

			if len(team1Info.MainRoster) == 0 || len(team2Info.MainRoster) == 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": "Both teams must have non-empty main rosters"})
				return
			}

			hasExistingCode, err := LeagueOfLegendsDAO.HasTournamentCode(getGameId(ctx))
			if checkErr(ctx, err) {
				return
			} else if hasExistingCode && !newCode {
				tournamentCode, err := LeagueOfLegendsDAO.GetTournamentCode(getGameId(ctx))
				if checkErr(ctx, err) {
					return
				} else {
					ctx.JSON(http.StatusOK, gin.H{"tournamentCode": tournamentCode})
				}
			} else {
				isRegistered, err := LeagueOfLegendsDAO.LeagueHasRegisteredTournament(getLeagueId(ctx))
				if checkErr(ctx, err) {
					return
				} else if !isRegistered {
					ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": "This league must first be registered as a league of legends tournament"})
					return
				}

				tournamentId, err := LeagueOfLegendsDAO.GetTournamentId(getLeagueId(ctx))
				if checkErr(ctx, err) {
					return
				}

				externalId := hex.EncodeToString(securecookie.GenerateRandomKey(32))
				fmt.Printf("adding external id: %v\n", externalId)
				if err = GameDAO.AddExternalId(getGameId(ctx), externalId); checkErr(ctx, err) {
					return
				}

				tournamentCode, err := LoLTournamentApi.CreateTournamentKey(tournamentId,
					fmt.Sprintf("{\"gameId\": \"%v\", \"team1Id\":%v, "+
						"\"team2Id\":%v,\"team1RefPlayerId\":\"%v\"}",
						externalId, gameInfo.Team1.TeamId, gameInfo.Team2.TeamId, team1Info.MainRoster[0].ExternalId))
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"errorDescription": err.Error()})
					return
				}

				if err = LeagueOfLegendsDAO.CreateTournamentCode(getGameId(ctx), tournamentCode); checkErr(ctx, err) {
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"tournamentCode": tournamentCode})
			}
		},
	}.createEndpointHandler()
}

func tournamentCallback() gin.HandlerFunc {
	//TODO: check if game is valid before doing anything else
	return func(ctx *gin.Context) {
		bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
		if checkErr(ctx, err) {
			return
		}
		if err = LoLTournamentApi.ForwardCompleteTournamentGame(bodyBytes); checkErr(ctx, err) {
			return
		} else {
			ctx.Status(http.StatusOK)
		}
	}
}

func RegisterLeagueOfLegendsHandlers(g *gin.RouterGroup) {
	g.POST("/registerTournament", registerTournament())
	g.POST("/teamsWithPlayers", createNewLoLTeamWithPlayers())
	g.POST("/receiveCompletedTournamentGame", receiveCompletedTournamentGame())
	g.POST("/tournamentCallback", tournamentCallback())
	g.GET("/stats/player", getPlayerStats)
	g.GET("/stats/team", getTeamStats)
	g.GET("/stats/champion", getChampionStats)
	g.POST("/games", createNewGame())
	g.GET("/games/:gameId/tournamentCode", storeGameId(), getTournamentCode())
	withTeamId := g.Group("/teams/:teamId", storeTeamId())
	withTeamId.POST("/players", leagueOfLegendsGetSummonerId(), createNewLoLPlayer())
	withTeamId.PUT("/players/:playerId", storePlayerId(), leagueOfLegendsGetSummonerId(), updateLoLPlayer())
	withTeamId.GET("/withRosters", getLoLTeamWithRosters())
}
