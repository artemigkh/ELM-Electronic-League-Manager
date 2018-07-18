package routes

import "github.com/gin-gonic/gin"

type GameInformation struct {
	Team1ID int `json:"team1Id"`
	Team2ID int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

func RegisterGameHandlers(g *gin.RouterGroup) {

}
