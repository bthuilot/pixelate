package matrix

import (
	"github.com/bthuilot/pixelate/vndr/rgbmatrix"
	"github.com/sirupsen/logrus"
	"image"
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
	//draw.Draw(c, c.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	//// don't forget  Render to agents the new led status
	//c.Render()
	s := &Service{
		Chan:   make(chan image.Image),
		Matrix: c,
	}
	go s.renderLoop()
	return s, nil
}

func (s *Service) ClearScreen() {
	s.Matrix.Clear()
}

func (s *Service) renderLoop() {
	//defer s.Matrix.Close()
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{PadLevelText: true})
	logger.WithField("service", "matrix")
	logger.Debug("beginning image loop")
	for {
		select {
		case drawImg := <-s.Chan:
			logger.Debug("image recieved, drawing")
			draw.Draw(s.Matrix, s.Matrix.Bounds(), drawImg, image.Point{}, draw.Src)
			if err := s.Matrix.Render(); err != nil {
				logger.Errorf("unable to render image: %s", err)
			}
		}
	}
}
