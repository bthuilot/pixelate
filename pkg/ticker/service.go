package ticker

import (
	"SpotifyDash/pkg/api"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
)

type Service struct {
	stock  string
	matrix chan image.Image
}

func (s *Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
	s.matrix = matrixChan
	s.stock = "BX"
	return nil
}

func (s *Service) Tick() (err error) {
	var img image.Image
	img, err = createImg(s.stock)
	if err != nil {
		return err
	}
	s.matrix <- img
	return nil
}

func createImg(text string) (image.Image, error) {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)

	if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 8); err != nil {
		return nil, err
	}

	dc.SetColor(color.White)
	dc.DrawStringAnchored(text, 0, 0, 0, 0) //maxWidth, 1.5, gg.AlignLeft)

	return dc.Image(), nil
}

func (s *Service) GetConfig() api.ConfigStore {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetConfig(config api.ConfigStore) error {
	//TODO implement me
	panic("implement me")
	return nil
}
