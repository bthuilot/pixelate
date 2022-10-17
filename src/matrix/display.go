package matrix

import (
	"github.com/bthuilot/pixelate/vndr/rgbmatrix"
)

// New creates a n  ew Display that renders to the matrix display
func NewDisplayCanvas() (*rgbmatrix.Canvas, error) {
	config := &rgbmatrix.DefaultConfig
	config.Cols = 64
	config.Rows = 64
	config.HardwareMapping = "adafruit-hat"
	config.Brightness = 50
	// create a new Matrix instance with the DefaultConfig
	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	return rgbmatrix.NewCanvas(matrix), err
}