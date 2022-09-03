// conductor contains all code for the Conductor interface, which is meant to server
// as a controller for starting, running and stopping renderers and the communication between
// the matrix service and displaying renderer
package conductor

import (
	"fmt"
	"log"

	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/bthuilot/pixelate/pkg/matrix"
)

// Conductor represents the controller for starting, stopping and running the
// renderers and matrix service.
type Conductor interface {
	// ListRenders will return the IDs of all available renderers
	ListRenders() []display.ID
	// InitNewRenderer will start a new renderer from its ID, and stop
	// any existing renderer
	InitNewRenderer(display.ID) error
	// GetCurrentRender will return the ID and config of the current renderer.
	// It also will return a boolean, indicating if the there is a renderer currently running.
	GetCurrentRenderer() (string, display.Config, bool)
	// UpdateConfig will update the config for the currently running renderer
	UpdateConfig(display.Config) error
	// GetSetup will return any links to put on a "setup" page incase the renderer needs additional configuration from the user
	// TODO(refator this)
	GetSetup() (display.SetupPage, bool)
	// StopCurrentRenderer will stop the currently running renderer
	StopCurrentRenderer() error
}

// conductor is the concrete implemenation of the Conductor interface.
// This seperation is done for testing purposes.
type conductor struct {
	// renderers is a mapping from Renderer ID to display.Renderer
	renderers map[string]display.Renderer
	// setup is the setup page for the currently running renderer
	setup display.SetupPage
	// matrix is a matrix.Service, which is reponsiblke for controlling the matrix screen
	matrix *matrix.Service
	// currentRenderer is the currently displaying renderer
	currentService *runningService
}

// runningService is a struct to represent the currently running services
type runningService struct {
	// channel is the chan to send display.Command to the running service,
	// which is used to control the service
	channel chan display.Command
	// id is the display.ID for the service
	id string
	// config is the current configuration for the renderer
	config display.Config
}

// SpawnConductor will construct a conductor wit hthe given matrix service and list
// of available renderers. Will fail if the IDs of the Renderers are not unique.
func SpawnConductor(mtrx *matrix.Service, rndrs []display.Renderer) Conductor {
	rndrMap := map[string]display.Renderer{}
	for _, s := range rndrs {
		name := s.GetName()
		if _, e := rndrMap[name]; e {
			log.Fatalf("service names should be unique, recieved %s", name)
		}
		rndrMap[name] = s
	}
	return conductor{
		renderers:      rndrMap,
		setup:          nil,
		currentService: nil,
		matrix:         mtrx,
	}
}

// GetSetup will return the current renderers display.SetupPage
func (c conductor) GetSetup() (setup display.SetupPage, running bool) {
	setup = c.setup
	running = c.currentService != nil
	return
}

// ListRenders returns a list of all display.IDs for all available renderers
func (c conductor) ListRenders() (result []display.ID) {
	for n := range c.renderers {
		result = append(result, n)
	}
	return
}

// InitNewRenderer will initialize a new renderer and stop any currently running
func (c conductor) InitNewRenderer(id display.ID) (err error) {
	svc, exist := c.renderers[id]
	if !exist {
		err = fmt.Errorf("cannot initialize no existant service %s", id)
	}
	_ = c.StopCurrentRenderer() // ignore the error, just want to stop a service if one is running
	c.setup = svc.Init(c.matrix.Chan)
	return nil
}

// GetCurrentRenderer will return the currently running renderer
func (c conductor) GetCurrentRenderer() (id string, config display.Config, isRunning bool) {
	if c.currentService != nil {
		isRunning = true
		id = c.currentService.id
		config = c.currentService.config
	}
	return
}

// UpdateConfig will update the configuration for the currently runnng renderer
func (c conductor) UpdateConfig(newCfg display.Config) (err error) {
	if c.currentService != nil {
		c.currentService.channel <- display.Command{
			Code:   display.Update,
			Config: newCfg,
		}
	} else {
		err = fmt.Errorf("no service running")
	}
	return
}

// StopCurrentRenderer will stop the currently running renderer, or return
// an error if none exists
func (c conductor) StopCurrentRenderer() (err error) {
	if c.currentService != nil {
		c.currentService.channel <- display.Command{
			Code: display.Stop,
		}
	} else {
		err = fmt.Errorf("no service is currently running")
	}
	return
}
