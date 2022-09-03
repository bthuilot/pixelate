package display

import (
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
	// GetName will return the ID for the Renderer
	GetName() ID
	Run(Config, chan Command, chan image.Image)
	GetDefaultConfig() Config
	GetTickInterval() time.Duration
	Init(chan image.Image) SetupPage
}

type ConfigType = int

type SetupPage []Attribute

type CommandCode = int

const (
	Stop CommandCode = iota
	Update
	Tick
)

type Command struct {
	Code   CommandCode
	Config Config
}
