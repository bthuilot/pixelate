package matrix

import (
	"SpotifyDash/internal/rgbmatrix"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

type Service struct {
	Chan   chan image.Image
	Matrix *rgbmatrix.Canvas
}

func CreateService() (*Service, error) {
	config := &rgbmatrix.DefaultConfig
	config.Cols = 64
	config.Rows = 64
	config.HardwareMapping = "adafruit-hat"
	// create a new Matrix instance with the DefaultConfig
	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return nil, err
	}

	// create the Canvas, implements the image.Image interface
	c := rgbmatrix.NewCanvas(m)

	// using the standard draw.Draw function we copy a white image onto the Canvas
	draw.Draw(c, c.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	// don't forget call Render to display the new led status
	c.Render()
	return &Service{
		Chan:   make(chan image.Image),
		Matrix: c,
	}, nil
}

func (s *Service) Init() {
	go func() {
		defer s.Matrix.Close()
		for {
			fmt.Println("Checking")
			select {
			case drawImg := <-s.Chan:
				fmt.Println("drawing")
				fmt.Println(drawImg.At(32, 32))
				draw.Draw(s.Matrix, s.Matrix.Bounds(), drawImg, image.Point{}, draw.Src)
				s.Matrix.Render()
			}
		}
	}()
}
