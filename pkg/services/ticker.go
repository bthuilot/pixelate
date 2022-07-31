package services

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/bthuilot/pixelate/pkg/util"

	"github.com/fogleman/gg"
)

type Ticker struct {
	matrix chan image.Image
	symbol string
	// TODO api
}

func (t Ticker) GetName() string {
	return "Stock Ticker"
}

func (t Ticker) Run(_ Config) chan Command {
	// TODO
	cmdChannel := make(chan Command)
	go func() {
	exit_routine:
		for {
			select {
			case cmd := <-cmdChannel:
				code := cmd.Code
				switch code {
				case Stop:
					break exit_routine
				case Tick:
					t.tick()
				case Update:
					// TODO get new ticker symbol
				}
			}
		}

	}()
	return cmdChannel
}

func (t Ticker) GetDefaultConfig() Config {
	return Config{
		"Ticker Symbol": "BX",
	}
}

func (t Ticker) Init(matrixChan chan image.Image) (page SetupPage) {
	// TODO init API client
	t.matrix = matrixChan
	return
}

func (t Ticker) tick() (err error) {
	ticker := "BX"
	if err != nil {
		img := util.RenderText("Please set a ticker")
		t.matrix <- img
		return err
	}
	if err != nil {
		img := util.RenderText("API Error")
		t.matrix <- img
		return err
	}
	t.matrix <- createImg(ticker, "-3", "100")
	return nil
}

func (t Ticker) GetTickInterval() time.Duration {
	return time.Minute * 10
}

/**************
 * Rendering *
 *************/

var (
	positiveChange = color.RGBA{G: 255, A: 255}
	negativeChange = color.RGBA{R: 255, A: 255}
)

var (
	symbolFont2, _ = gg.LoadFontFace(util.BankPrinterFontPath, 32)
	symbolFont4, _ = gg.LoadFontFace(util.BankPrinterFontPath, 16)
	changeFont, _  = gg.LoadFontFace(util.BankPrinterFontPath, 10)
	priceFont, _   = gg.LoadFontFace(util.BankPrinterFontPath, 15)
)

func createImg(ticker string, change string, price string) image.Image {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)
	font := symbolFont2
	if len(ticker) > 2 {
		font = symbolFont4
	}
	dc.SetFontFace(font)
	dc.SetColor(color.White)
	// Ticker
	tW, tH := dc.MeasureString(ticker)
	tH -= (tH / 5)
	dc.DrawString(ticker, 0, tH)
	// Price
	dc.SetFontFace(priceFont)
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
	dc.SetFontFace(changeFont)
	cW, cH := dc.MeasureString(change)
	dc.DrawStringAnchored(change, 64-cW, pY-(cH/2.0+pH/2.0), 0, 0)
	// Arrow
	radius := math.Min((64-cW)/2.0, (cH / 2.0))
	dc.DrawRegularPolygon(3, tW+radius, (tH*0.95)/5.0+radius/2, radius, rotation)
	dc.Fill()
	dc.Stroke()

	return dc.Image()
}

func renderError(msg string) image.Image {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)
	dc.SetFontFace(util.ErrorFont)
	dc.SetColor(color.White)
	dc.DrawStringWrapped(msg, 32, 32, 0.5, 0.5, 64, 1.0, gg.AlignCenter)
	return dc.Image()
}
