package ticker

import (
	"SpotifyDash/pkg/util"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math"
)

var (
	SymbolFont2, _ = gg.LoadFontFace(util.BankPrinterFontPath, 32)
	SymbolFont4, _ = gg.LoadFontFace(util.BankPrinterFontPath, 16)
	ChangeFont, _  = gg.LoadFontFace(util.BankPrinterFontPath, 10)
	PriceFont, _   = gg.LoadFontFace(util.BankPrinterFontPath, 15)
)

func createImg(ticker string, change string, price string) image.Image {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)
	font := SymbolFont2
	if len(ticker) > 2 {
		font = SymbolFont4
	}
	dc.SetFontFace(font)
	dc.SetColor(color.White)
	// Ticker
	tW, tH := dc.MeasureString(ticker)
	tH -= (tH / 5)
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
