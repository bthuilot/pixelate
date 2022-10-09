package util

import (
	"github.com/fogleman/gg"
	"path"
)

var FontDir = path.Join(GetDir(), "assets/fonts")

var BankPrinterFontPath = FontDir + "/F25_Bank_Printer.ttf"

var ErrorFont, _ = gg.LoadFontFace(BankPrinterFontPath, 12)
