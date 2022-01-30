package image_util

import (
	"github.com/disintegration/imaging"
	"image"
)

func Resize(img image.Image) image.Image {
	return imaging.Resize(img, 64, 64, imaging.Lanczos)
}