// conductor contains all code for the Conductor interface, which is meant to server
// as a controller for starting, running and stopping renderers and the communication between
// the matrix service and displaying renderer
package conductor

import (
	"fmt"
	"log"
	"time"

	"github.com/bthuilot/pixelate/matrix"
	"github.com/bthuilot/pixelate/rendering"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Conductor represents the controller for starting, stopping and running the
// renderers and matrix service.
type Conductor interface {
	// ListRenderer will return the IDs of all available renderers
	ListRenderers() []rendering.ID
	// InitNewRenderer will start a new renderer from its ID, and stop
	// any existing renderer
	InitNewRenderer(rendering.ID) error
	// GetCurrentRenderer will return the ID and config of the current renderer.
	// It also will return a boolean, indicating if the there is a renderer currently running.
	GetCurrentRenderer() (string, rendering.Config, []rendering.ConfigAttribute, bool)
	// UpdateConfig will update the config for the currently running renderer
	UpdateConfig(rendering.Config) error
	// StopCurrentRenderer will stop the currently running renderer
	StopCurrentRenderer() error

	RegisterAgentEndpoints(r *gin.Engine)
}

// conductor is the concrete implementation of the Conductor interface.
// This separation is done for testing purposes.
type conductor struct {
	// renderers is a mapping from Agent ID to rendering.Agent
	renderers map[string]rendering.Agent
	// matrix is a matrix.Service, which is responsible for controlling the matrix screen
	matrix *matrix.Service
	// currentRenderer is the currently displaying renderer
	currentAgent *currentAgent
}

type currentAgent struct {
	renderer rendering.Agent
	exitChan chan interface{}
}

// New will construct a conductor wit hthe given matrix service and list
// of available renderers. Will fail if the IDs of the Renderers are not unique.
func New(mtrx *matrix.Service, rndrs []rendering.Agent) Conductor {
	rndrMap := map[string]rendering.Agent{}
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
func (c *conductor) ListRenderers() (result []rendering.ID) {
	for n := range c.renderers {
		result = append(result, n)
	}
	return
}

// InitNewRenderer will initialize a new renderer and stop any currently running
func (c *conductor) InitNewRenderer(id rendering.ID) error {
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

func (c *conductor) rendererLoop(agent rendering.Agent, exitChan chan interface{}) {
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
func (c *conductor) GetCurrentRenderer() (id string, config rendering.Config, attrs []rendering.ConfigAttribute, isRunning bool) {
	if isRunning = c.currentAgent != nil; isRunning {
		id = c.currentAgent.renderer.GetName()
		config = c.currentAgent.renderer.GetConfig()
		attrs = c.currentAgent.renderer.GetAdditionalConfig()
	}
	return
}

// UpdateConfig will update the configuration for the currently runnnig renderer
func (c *conductor) UpdateConfig(newCfg rendering.Config) error {
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
