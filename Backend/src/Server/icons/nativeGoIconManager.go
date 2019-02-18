package icons

import (
	"Server/config"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image/png"
	"net/http"
	"os"
)

type GoIconManager struct {
	OutPath string
}

func CreateGoIconManager(conf config.Config) IconManager {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", os.ModePerm)
	}

	return &GoIconManager{
		OutPath: conf.GetIconsDir(),
	}
}

func (m *GoIconManager) StoreNewIcon(ctx *gin.Context) (string, string, error) {
	tmpFileLoc := "tmp/" + randomdata.RandStringRunes(10)
	formFile, err := ctx.FormFile("icon")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return "", "", err
	}

	if err := ctx.SaveUploadedFile(formFile, tmpFileLoc); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return "", "", err
	}

	defer os.Remove(tmpFileLoc)

	file, err := os.Open(tmpFileLoc)
	if err != nil {
		return "", "", err
	}

	img, err := png.Decode(file)
	if err != nil {
		return "", "", err
	}
	file.Close()

	largeOutPath := randomdata.RandStringRunes(10) + ".png"
	smallOutPath := randomdata.RandStringRunes(10) + ".png"

	largeImg := resize.Resize(512, 512, img, resize.Lanczos3)
	smallImg := resize.Resize(128, 128, img, resize.Lanczos3)

	largeOut, err := os.Create(m.OutPath + largeOutPath)

	if err != nil {
		return "", "", err
	}
	defer largeOut.Close()

	smallOut, err := os.Create(m.OutPath + smallOutPath)
	if err != nil {
		return "", "", err
	}
	defer smallOut.Close()

	png.Encode(largeOut, largeImg)
	png.Encode(smallOut, smallImg)

	return smallOutPath, largeOutPath, nil
}
