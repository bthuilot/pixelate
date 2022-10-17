// Package conductor contains all code for the Conductor interface, which is meant to server
// as a controller for starting, running and stopping renderers and the communication between
// the display service and displaying agent
package matrix

import (
	"fmt"
	"github.com/bthuilot/pixelate/vndr/rgbmatrix"
	"image"
	"image/draw"
	"time"

	"github.com/bthuilot/pixelate/rendering"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Conductor represents the controller for starting, stopping and running the
// renderers and display service.
type Conductor interface {
	// ListAgents will return the IDs of all available renderers
	ListAgents() []rendering.ID
	// InitNewAgent will start a new agent from its ID, and stop
	// any existing agent
	InitNewAgent(rendering.ID) error
	// GetCurrentAgent will return the ID and config of the current agent.
	// It also will return a boolean, indicating if the there is a agent currently running.
	GetCurrentAgent() (string, rendering.Config, []rendering.ConfigAttribute, bool)
	// UpdateConfig will update the config for the currently running agent
	UpdateConfig(rendering.Config) error
	// StopCurrentAgent  will stop the currently running agent
	StopCurrentAgent() error
	// RegisterAgentEndpoints will register custom agent specific endpoints against
	// the HTTP server
	RegisterAgentEndpoints(r *gin.Engine)
}

// conductor is the concrete implementation of the Conductor interface.
// This separation is done for testing purposes.
type conductor struct {
	// renderers is a mapping from Agent ID to rendering.Agent
	renderers map[string]rendering.Agent
	// display is a matrix.Display, which is responsible for controlling the display screen
	display *rgbmatrix.Canvas
	// currentRenderer is the currently displaying agent
	currentAgent *currentAgent
}

// currentAgent represents the currently running agent
// and the channel used to communicate when the agent should be stopped
type currentAgent struct {
	// agent is currently running agent
	agent rendering.Agent
	// exitChan is the channel used to communicate to the rendering loop to stop
	exitChan chan interface{}
}

// NewConductor will construct a conductor wit hthe given display service and list
// of available renderers. Will fail if the IDs of the Renderers are not unique.
func NewConductor(rndrs []rendering.Agent) (Conductor, error) {
	rndrMap := map[string]rendering.Agent{}
	for _, s := range rndrs {
		name := s.GetName()
		if _, e := rndrMap[name]; e {
			return nil, fmt.Errorf("service names should be unique, recieved %s", name)
		}
		rndrMap[name] = s
	}
	canvas, err := NewDisplayCanvas()
	return &conductor{
		renderers:    rndrMap,
		currentAgent: nil,
		display:      canvas,
	}, err
}

// RegisterAgentEndpoints will call RegisterEndpoints on every Agent,
// allowing them to  register endpoints against the HTTP server
func (c *conductor) RegisterAgentEndpoints(r *gin.Engine) {
	logrus.Info("registering agent endpoints ")
	for _, a := range c.renderers {
		logrus.Infof("registering endpoints for %s", a.GetName())
		a.RegisterEndpoints(r)
	}
}

// ListAgents returns a list of all agents.IDs for all available renderers
func (c *conductor) ListAgents() (result []rendering.ID) {
	for n := range c.renderers {
		result = append(result, n)
	}
	return
}

// InitNewAgent will initialize a new rendering agent and stop any currently running
func (c *conductor) InitNewAgent(id rendering.ID) error {
	logrus.Infof("starting rendering agent %s", id)
	agent, exist := c.renderers[id]
	if !exist {
		return fmt.Errorf("initialize non existant service %s", id)
	}
	_ = c.StopCurrentAgent() // ignore the error, just want to stop a service if one is running
	exit := make(chan interface{})
	c.currentAgent = &currentAgent{agent: agent, exitChan: exit}
	go c.rendererLoop(agent, exit)
	return nil
}

// rendererLoop will begin a loop for continually drawing frames to the display
func (c *conductor) rendererLoop(agent rendering.Agent, exitChan chan interface{}) {
	logrus.Infof("beginning polling loop for %s", agent.GetName())
	for {
		select {
		case <-exitChan:
			logrus.Info("stopping agent %s\n", agent.GetName())
			return
		default:
			logrus.Debug("performing render tick")
			frame := agent.NextFrame()
			draw.Draw(c.display, c.display.Bounds(), frame, image.Point{}, draw.Src)
			if err := c.display.Render(); err != nil {
				logrus.Errorf("unable to render image: %s", err)
			}
		}
		logrus.Debug("sleeping for %s", agent.GetTick())
		time.Sleep(agent.GetTick())
	}
}

// GetCurrentAgent will return the currently running agent
func (c *conductor) GetCurrentAgent() (id string, config rendering.Config, attrs []rendering.ConfigAttribute, isRunning bool) {
	if isRunning = c.currentAgent != nil; isRunning {
		id = c.currentAgent.agent.GetName()
		config = c.currentAgent.agent.GetConfig()
		attrs = c.currentAgent.agent.GetAdditionalConfig()
	}
	return
}

// UpdateConfig will update the configuration for the currently runnnig agent
func (c *conductor) UpdateConfig(newCfg rendering.Config) error {
	if c.currentAgent == nil {
		return fmt.Errorf("no service running")
	}
	return c.currentAgent.agent.SetConfig(newCfg)
}

// StopCurrentAgent will stop the currently running agent, or return
// an error if none exists
func (c *conductor) StopCurrentAgent() error {
	if c.currentAgent == nil {
		return fmt.Errorf("no service is currently running")
	}
	c.currentAgent.exitChan <- struct{}{}
	c.currentAgent = nil
	return nil
}
