package icons

import (
	"github.com/gin-gonic/gin"
	"image"
)

type IconManager interface {
	StoreNewIconFromBase64String(icon string) (string, string, error)
	StoreNewIconFromForm(ctx *gin.Context) (string, string, error)
	StoreNewIconFromImage(img image.Image) (string, string, error)
}
