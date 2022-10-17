package rendering

import (
	"embed"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"path"
)

//go:embed fonts/*
var fontFiles embed.FS

const FontDir = "fonts"

var ErrorFont, _ = LoadFont("F25_Bank_Printer.ttf", 12)

func LoadFont(name string, points float64) (font.Face, error) {
	fontBytes, err := fontFiles.ReadFile(path.Join(FontDir, name))
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
