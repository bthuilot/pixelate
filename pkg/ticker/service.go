package ticker

import (
	"SpotifyDash/pkg/api"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"math"
	"time"
)

type Service struct {
	stock  string
	matrix chan image.Image
}

var (
	positiveChange = color.RGBA{G: 255, A: 255}
	negativeChange = color.RGBA{R: 255, A: 255}
)

const fontURI = "/home/bryce/github/PiMatrix/assets/fonts/F25_Bank_Printer.ttf"

var (
	SymbolFont, _ = gg.LoadFontFace(fontURI, 32)
	ChangeFont, _ = gg.LoadFontFace(fontURI, 10)
	PriceFont, _  = gg.LoadFontFace(fontURI, 15)
)

func (s *Service) Init(matrixChan chan image.Image, engine *gin.Engine) error {
	s.matrix = matrixChan
	s.stock = "BX"
	return nil
}

func (s *Service) Tick() (err error) {
	var img image.Image
	info, err := getStockInfo(s.stock)
	if err != nil {
		return err
	}
	fmt.Println(info)
	img, err = createImg(s.stock, info.Change, info.Price[:len(info.Price)-2])
	if err != nil {
		return err
	}
	s.matrix <- img
	return nil
}

func createImg(ticker string, change string, price string) (image.Image, error) {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)
	dc.SetFontFace(SymbolFont)
	dc.SetColor(color.White)
	// Ticker
	tW, tH := dc.MeasureString(ticker)
	tH -= (tH / 5)
	fmt.Printf("%f\n", tH/2.0)
	dc.DrawString(ticker, 0, tH)
	// Price
	dc.SetFontFace(PriceFont)
	pW, pH := dc.MeasureString(price)
	pH -= (pH / 5.0)
	pX, pY := 64-(pW), 64-(pH/3.0)
	dc.DrawString(price, pX, pY+2)

	// Change //
	changeColor, rotation := negativeChange, math.Pi
	if len(change) > 0 && change[0] != '-' {
		change = "+" + change
		changeColor = positiveChange
		rotation = 0
	}
	dc.SetColor(changeColor)
	// Amount
	dc.SetFontFace(ChangeFont)
	fmt.Println("Test")
	fmt.Println(dc.FontHeight())
	cW, cH := dc.MeasureString(change)
	dc.DrawStringAnchored(change, 64-cW, pY-(cH/2.0+pH/2.0), 0, 0)
	// Arrow
	radius := math.Min((64-cW)/2.0, (cH / 2.0))
	dc.DrawRegularPolygon(3, tW+radius, (tH*0.95)/5.0+radius/2, radius, rotation)
	dc.Fill()
	dc.Stroke()

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

func (s *Service) RefreshDelay() time.Duration {
	return time.Minute * 5
}
