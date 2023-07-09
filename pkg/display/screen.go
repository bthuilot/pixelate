package display

import (
	"image"
	"time"

	"github.com/gin-gonic/gin"
)

// Screen represents a service that produces an [image.Image] to
// be rendered to the LED matrix
type Screen interface {
	// Init allows the screen to initialize and register custom
	// HTTP routes in the given [gin.RouterGroup]. This method
	// should return a string that represents the formatted name
	// for this screen
	Init(r *gin.RouterGroup) (string, error)
	// Render will produce an [image.Image] to be displayed onto
	// the LED matrix. the [time.Duration] returned represents how
	// long the display should wait before calling again to update.
	// Values of [time.Duration] that are 0 or less will result in no
	// more update calls being made
	Render() (image.Image, time.Duration, error)
	// GetConfig returns the current configuration of the Screen
	GetConfig() map[string]string
	// UpdateConfig will update the configuration for the Screen
	UpdateConfig(map[string]string) error
	// GetHTMLPage returns the HTMLAtrributes that should
	// be rendered on the home page for control of the Screen
	GetHTMLPage() []HTMLAttribute
}
