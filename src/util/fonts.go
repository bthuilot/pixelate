package util

import (
	"path"

	"github.com/fogleman/gg"
)

var FontDir = path.Join(GetDir(), "assets/fonts")

var BankPrinterFontPath = FontDir + "/F25_Bank_Printer.ttf"

var ErrorFont, _ = gg.LoadFontFace(BankPrinterFontPath, 12)
