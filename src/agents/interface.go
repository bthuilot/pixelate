package agents

import (
	"github.com/gin-gonic/gin"
	"image"
	"time"
)

// ID is an identifier used for labeling renderers
type ID = string

// Config is a configraution for a Renderer, which is a mapping
// of config keys to values
type Config map[string]string

// Renderer represents a particular service to draw to the board.
// Renderers will start a go rountine a use the provided channels to communicate with the
type Renderer interface {
	// GetName will return the ID for the
	GetName() ID
	GetConfig() Config
	SetConfig(Config) error
	GetAdditionalHTML() []Attribute
	Render(chan image.Image)
	GetTick() time.Duration
	RegisterEndpoints(r *gin.Engine)
}
