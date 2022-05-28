package util

import "github.com/fogleman/gg"

const FontDir = "/home/bryce/github/PiMatrix/assets/fonts"

const BankPrinterFontPath = FontDir + "/F25_Bank_Printer.ttf"

var ErrorFont, _ = gg.LoadFontFace(BankPrinterFontPath, 12)
