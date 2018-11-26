package icons

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/nfnt/resize"
	"image/png"
	"os"
)

type GoIconManager struct {
	OutPath string
}

func (m *GoIconManager) StoreNewIcon(tempLoc string) (string, string, error) {
	defer os.Remove(tempLoc)

	file, err := os.Open(tempLoc)
	if err != nil {
		return "", "", err
	}

	img, err := png.Decode(file)
	if err != nil {
		return "", "", err
	}
	file.Close()

	largeOutPath := randomdata.RandStringRunes(10)
	smallOutPath := randomdata.RandStringRunes(10)

	largeImg := resize.Resize(512, 512, img, resize.Lanczos3)
	smallImg := resize.Resize(128, 128, img, resize.Lanczos3)

	largeOut, err := os.Create(m.OutPath + largeOutPath + ".png")
	if err != nil {
		return "", "", err
	}
	defer largeOut.Close()

	smallOut, err := os.Create(m.OutPath + smallOutPath + ".png")
	if err != nil {
		return "", "", err
	}
	defer smallOut.Close()

	png.Encode(largeOut, largeImg)
	png.Encode(smallOut, smallImg)

	return largeOutPath, smallOutPath, nil
}
