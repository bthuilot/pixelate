package ticker

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"time"
)

type Service struct {
	matrix chan image.Image
	config api.ConfigStore
}

var (
	positiveChange = color.RGBA{G: 255, A: 255}
	negativeChange = color.RGBA{R: 255, A: 255}
)

func (s *Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
	s.matrix = matrixChan
	s.config = api.ConfigStore{
		"ticker": "BX",
	}
	return nil
}

func (s *Service) getStockTicker() (string, error) {
	if ticker, c := s.config["ticker"]; c {
		return ticker, nil
	}
	return "", fmt.Errorf("invalid ticker")
}

func (s *Service) Tick() (err error) {
	ticker, err := s.getStockTicker()
	if err != nil {
		img := util.RenderText("Please set a ticker")
		s.matrix <- img
		return err
	}
	info, err := getStockInfo(ticker)
	if err != nil {
		img := util.RenderText("API Error")
		s.matrix <- img
		return err
	}
	s.matrix <- createImg(ticker, info.Change, info.Price[:len(info.Price)-2])
	return nil
}

func (s *Service) GetConfig() api.ConfigStore {
	return s.config
}

func (s *Service) SetConfig(config api.ConfigStore) error {
	s.config = config
	return nil
}

func (s *Service) RefreshDelay() time.Duration {
	return time.Minute * 5
}

func (s Service) GetID() string {
	return "ticker"
}
