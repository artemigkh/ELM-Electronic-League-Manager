package routesTest

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

type idResponse struct {
	Id int `json:"id"`
}

var router *gin.Engine
