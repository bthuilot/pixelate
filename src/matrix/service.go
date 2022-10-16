package matrix

import (
	"github.com/bthuilot/pixelate/vndr/rgbmatrix"
	"github.com/sirupsen/logrus"
	"image"
	"image/draw"
)

// Service represents the running service that interacts with the matrix display
type Service struct {
	// Chan is the channel used to receive images to render to the display
	Chan chan image.Image
	// Matrix is the actual canvas of the display to draw to
	Matrix *rgbmatrix.Canvas
	// Exit is the channel used to communicate that the rendering Loop should be exited
	Exit chan interface{}
}

// New creates a new Service that renders to the matrix display
func New() (*Service, error) {
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
	s := &Service{
		Chan:   make(chan image.Image),
		Matrix: rgbmatrix.NewCanvas(m),
		Exit:   make(chan interface{}),
	}
	go s.renderLoop()
	return s, nil
}

func (s *Service) ClearScreen() {
	s.Matrix.Clear()
}

func (s *Service) renderLoop() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{PadLevelText: true})
	logger.WithField("service", "matrix")
	logger.Debug("beginning image loop")
	for {
		select {
		case <-s.Exit:
			logger.Warning("Received exit, shutting down gracefully")
			if err := s.Matrix.Clear(); err != nil {
				logger.Error(err)
			}
		case drawImg := <-s.Chan:
			logger.Debug("image received, drawing")
			draw.Draw(s.Matrix, s.Matrix.Bounds(), drawImg, image.Point{}, draw.Src)
			if err := s.Matrix.Render(); err != nil {
				logger.Errorf("unable to render image: %s", err)
			}
		}
	}
}
