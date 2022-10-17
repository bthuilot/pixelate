package rendering

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"image"
	"time"
)

// ID is an identifier used for labeling agents
type ID = string

// Config is a configuration for an Agent, which is a mapping
// of config keys to values
type Config map[string]string

// Agent represents a service taht will draw to the matrix
type Agent interface {
	// GetName will return the ID for the
	GetName() ID
	// GetConfig will return the current configuration of the Agent
	GetConfig() Config
	// SetConfig will set the configuration of the agent
	SetConfig(Config) error
	// GetAdditionalConfig will return a list of additional ConfigAttributes to display on the
	// config page
	GetAdditionalConfig() []ConfigAttribute
	// NextFrame will return the next image.Image to render to the display
	NextFrame() image.Image
	// GetTick will return the duration to sleep for between drawings
	GetTick() time.Duration
	// RegisterEndpoints will register agent specific endpoints to the HTTP Server
	RegisterEndpoints(r *gin.Engine)
}

// ConfigAttribute is an additional HTML attribute to display on the
// config page
type ConfigAttribute interface {
	// GetHTML will return the html to be displayed
	GetHTML() string
}

// Link is an HTML Link
type Link struct {
	// Name is the text to display as the link
	Name string
	// Href is the location of the link
	Href string
}

func (b Link) GetHTML() string {
	return fmt.Sprintf("<a href='%s'>%s</a>",
		html.EscapeString(b.Href), html.EscapeString(b.Name))
}

// Text represents plain text to be displayed
type Text struct {
	// Content is the content of the text blob
	Content string
}

func (t Text) GetHTML() string {
	return html.EscapeString(t.Content)
}
