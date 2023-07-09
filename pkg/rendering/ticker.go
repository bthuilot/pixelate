package rendering

import (
	"fmt"
	"github.com/bthuilot/pixelate/pkg/alphavantage"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/fogleman/gg"
)

type Ticker struct {
	config Config
}

const tickerSymbolConfigName = "TickerSymbol"

func NewTickerAgent() Agent {
	return &Ticker{
		Config{
			tickerSymbolConfigName: "BX",
		},
	}
}

const TickerAgentID ID = "StockTicker"

func (t *Ticker) GetName() ID {
	return TickerAgentID
}

func (t *Ticker) GetConfig() Config {
	return t.config
}

func (t *Ticker) SetConfig(config Config) error {
	if _, exist := config[tickerSymbolConfigName]; !exist {
		return fmt.Errorf("invalid configuration, must contain '%s'", tickerSymbolConfigName)
	}
	t.config = config
	return nil
}

func (t *Ticker) GetAdditionalConfig() (attrs []ConfigAttribute) {
	return
}

func (t *Ticker) NextFrame() image.Image {
	var ticker string
	var exist bool
	if ticker, exist = t.config[tickerSymbolConfigName]; !exist {
		return RenderText("Please set a ticker")
	}
	res, err := alphavantage.GetStockInfo(ticker)
	if err != nil {
		return RenderText("API Error")
	}
	return createImg(ticker, res.Change, res.Price)
}

func (t *Ticker) GetTick() time.Duration {
	return time.Minute * 10
}

func (t *Ticker) RegisterEndpoints(_ *gin.Engine) {
	// No endpoints needed
}

/**************
 * Rendering *
 *************/

var (
	positiveChange = color.RGBA{G: 255, A: 255}
	negativeChange = color.RGBA{R: 255, A: 255}
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
	dc.SetFontFace(ErrorFont)
	dc.SetColor(color.White)
	dc.DrawStringWrapped(msg, 32, 32, 0.5, 0.5, 64, 1.0, gg.AlignCenter)
	return dc.Image()
}
