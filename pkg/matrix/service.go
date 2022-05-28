package matrix

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/internal/rgbmatrix"
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
	config.Brightness = 50
	// create a new Matrix instance with the DefaultConfig
	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return nil, err
	}

	// create the Canvas, implements the image.Image interface
	c := rgbmatrix.NewCanvas(m)

	// using the standard draw.Draw function we copy a white image onto the Canvas
	draw.Draw(c, c.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	// don't forget call Render to display the new led status
	c.Render()
	return &Service{
		Chan:   make(chan image.Image),
		Matrix: c,
	}, nil
}

func (s *Service) ClearScreen() {
	s.Matrix.Clear()
}

func (s *Service) Init() {
	go func() {
		defer s.Matrix.Close()
		for {
			logging.InfoLogger.Println("Polling for img")
			select {
			case drawImg := <-s.Chan:
				logging.InfoLogger.Println("Image received, drawing")
				draw.Draw(s.Matrix, s.Matrix.Bounds(), drawImg, image.Point{}, draw.Src)
				s.Matrix.Render()
			}
		}
	}()
}
