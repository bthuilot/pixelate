// conductor contains all code for the Conductor interface, which is meant to server
// as a controller for starting, running and stopping renderers and the communication between
// the matrix service and displaying renderer
package conductor

import (
	"fmt"
	"github.com/bthuilot/pixelate/pkg/agents"
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

// Conductor represents the controller for starting, stopping and running the
// renderers and matrix service.
type Conductor interface {
	// ListRenderer will return the IDs of all available renderers
	ListRenderers() []agents.ID
	// InitNewRenderer will start a new renderer from its ID, and stop
	// any existing renderer
	InitNewRenderer(agents.ID) error
	// GetCurrentRenderer will return the ID and config of the current renderer.
	// It also will return a boolean, indicating if the there is a renderer currently running.
	GetCurrentRenderer() (string, agents.Config, []agents.Attribute, bool)
	// UpdateConfig will update the config for the currently running renderer
	UpdateConfig(agents.Config) error
	// StopCurrentRenderer will stop the currently running renderer
	StopCurrentRenderer() error

	RegisterAgentEndpoints(r *gin.Engine)
}

// conductor is the concrete implementation of the Conductor interface.
// This separation is done for testing purposes.
type conductor struct {
	// renderers is a mapping from Renderer ID to agents.Renderer
	renderers map[string]agents.Renderer
	// matrix is a matrix.Service, which is responsible for controlling the matrix screen
	matrix *matrix.Service
	// currentRenderer is the currently displaying renderer
	currentAgent *currentAgent
}

type currentAgent struct {
	renderer agents.Renderer
	exitChan chan interface{}
}

// SpawnConductor will construct a conductor wit hthe given matrix service and list
// of available renderers. Will fail if the IDs of the Renderers are not unique.
func SpawnConductor(mtrx *matrix.Service, rndrs []agents.Renderer) Conductor {
	rndrMap := map[string]agents.Renderer{}
	for _, s := range rndrs {
		name := s.GetName()
		if _, e := rndrMap[name]; e {
			log.Fatalf("service names should be unique, recieved %s", name)
		}
		rndrMap[name] = s
	}
	return &conductor{
		renderers:    rndrMap,
		currentAgent: nil,
		matrix:       mtrx,
	}
}

func (c *conductor) RegisterAgentEndpoints(r *gin.Engine) {
	logrus.Info("registering agent endpoints ")
	for _, s := range c.renderers {
		logrus.Infof("registering endpoints for %s", s.GetName())
		s.RegisterEndpoints(r)
	}
}

// ListRenderers returns a list of all agents.IDs for all available renderers
func (c *conductor) ListRenderers() (result []agents.ID) {
	for n := range c.renderers {
		result = append(result, n)
	}
	return
}

// InitNewRenderer will initialize a new renderer and stop any currently running
func (c *conductor) InitNewRenderer(id agents.ID) error {
	logrus.Infof("starting rendering agent %s", id)
	agent, exist := c.renderers[id]
	if !exist {
		return fmt.Errorf("initialize non existant service %s", id)
	}
	_ = c.StopCurrentRenderer() // ignore the error, just want to stop a service if one is running
	exit := make(chan interface{})
	c.currentAgent = &currentAgent{renderer: agent, exitChan: exit}
	go c.rendererLoop(agent, exit)
	return nil
}

func (c *conductor) rendererLoop(agent agents.Renderer, exitChan chan interface{}) {
	logrus.Infof("beginning polling loop for %s", agent.GetName())
	for {
		select {
		case <-exitChan:
			logrus.Info("stopping agent %s\n", agent.GetName())
			return
		default:
			logrus.Debug("performing render tick")
			agent.Render(c.matrix.Chan)
		}
		logrus.Debug("slepping for %s", agent.GetTick())
		time.Sleep(agent.GetTick())
	}
}

// GetCurrentRenderer will return the currently running renderer
func (c *conductor) GetCurrentRenderer() (id string, config agents.Config, attrs []agents.Attribute, isRunning bool) {
	if isRunning = c.currentAgent != nil; isRunning {
		id = c.currentAgent.renderer.GetName()
		config = c.currentAgent.renderer.GetConfig()
		attrs = c.currentAgent.renderer.GetAdditionalHTML()
	}
	return
}

// UpdateConfig will update the configuration for the currently runnnig renderer
func (c *conductor) UpdateConfig(newCfg agents.Config) error {
	if c.currentAgent == nil {
		return fmt.Errorf("no service running")
	}
	return c.currentAgent.renderer.SetConfig(newCfg)
}

// StopCurrentRenderer will stop the currently running renderer, or return
// an error if none exists
func (c *conductor) StopCurrentRenderer() error {
	if c.currentAgent == nil {
		return fmt.Errorf("no service is currently running")
	}
	c.currentAgent.exitChan <- struct{}{}
	c.currentAgent = nil
	return nil
}
