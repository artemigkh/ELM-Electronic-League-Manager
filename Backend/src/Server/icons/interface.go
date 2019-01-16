package icons

import "github.com/gin-gonic/gin"

type IconManager interface {
	StoreNewIcon(ctx *gin.Context) (string, string, error)
}
