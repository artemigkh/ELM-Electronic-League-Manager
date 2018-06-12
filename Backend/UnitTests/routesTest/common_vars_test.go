package routesTest

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

var router *gin.Engine
