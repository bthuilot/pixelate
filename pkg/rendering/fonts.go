package rendering

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io"
	"io/fs"
)

const BankPrinterFontName = "F25_Bank_Printer.ttf"

var (
	ErrorFont   font.Face
	symbolFont2 font.Face
	symbolFont4 font.Face
	changeFont  font.Face
	priceFont   font.Face
)

var font1 *truetype.Font

func LoadFonts(fonts fs.FS) error {
	fontFile, err := fonts.Open(BankPrinterFontName)
	if err != nil {
		return err
	}
	defer fontFile.Close()

	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		return err
	}
	font1, err = truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	symbolFont2, _ = LoadFont(32)
	symbolFont4, _ = LoadFont(16)
	changeFont, _ = LoadFont(10)
	priceFont, _ = LoadFont(15)
	if ErrorFont, err = LoadFont(12); err != nil {
		return err
	}
	return err
}

func LoadFont(points float64) (font.Face, error) {
	face := truetype.NewFace(font1, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
