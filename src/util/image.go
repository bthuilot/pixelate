package util

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"net/http"
)

func Resize(img image.Image) image.Image {
	return imaging.Resize(img, 64, 64, imaging.Lanczos)
}

func FromURL(url string) (image.Image, error) {
	//Get the response bytes from the url
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("unable to close response")
		}
	}(response.Body)

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}
	img, _, err := image.Decode(response.Body)
	return img, err
}

func RenderText(text string) image.Image {
	dc := gg.NewContext(64, 64)
	dc.DrawImage(&image.Uniform{C: color.Black}, 0, 0)
	dc.SetFontFace(ErrorFont)
	dc.SetColor(color.White)
	dc.DrawStringWrapped(text, 32, 32, 0.5, 0.5, 64, 1.0, gg.AlignCenter)
	return dc.Image()
}
