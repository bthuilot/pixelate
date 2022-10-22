package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
	"net/http"
)

// registerEndpoints will register all HTTP endpoints for the server
func (s Server) registerEndpoints(staticDir fs.FS) {
	logrus.Info("registering HTTP endpoints")
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	s.router.StaticFS("/static", http.FS(staticDir))
	s.router.GET("/", s.RenderDashboard)

	// Agents //
	// All Agents
	s.router.GET("/agents", s.ListAgents)
	// Current Agent
	s.router.GET("/agents/current", s.GetCurrentAgent)
	s.router.DELETE("/agents/current", s.StopCurrentAgent)
	s.router.POST("/agents/current", s.SetAgent)
	// Config
	s.router.POST("/agents/current/config", s.UpdateConfig)
}

// RenderDashboard is the endpoint to load and render the dashboard template
func (s Server) RenderDashboard(c *gin.Context) {
	name, cfg, attrs, running := s.cndtr.GetCurrentAgent()
	var html []template.HTML
	for _, attr := range attrs {
		html = append(html, template.HTML(attr.GetHTML()))
	}
	c.HTML(http.StatusOK, "index.tmpl", struct {
		CurrentAgentRunning bool
		CurrentAgent        string
		Config              map[string]string
		Attributes          []template.HTML
		Agents              []string
	}{
		CurrentAgentRunning: running,
		CurrentAgent:        name,
		Config:              cfg,
		Attributes:          html,
		Agents:              s.cndtr.ListAgents(),
	})
}

// GetCurrentAgent is the endpoint to return the currently running agent
func (s Server) GetCurrentAgent(c *gin.Context) {
	id, cfg, _, isRunning := s.cndtr.GetCurrentAgent()
	c.JSON(http.StatusOK, ValidResponse[CurrentAgentResponse]{
		Success: true,
		Response: CurrentAgentResponse{
			ID:        id,
			Config:    cfg,
			IsRunning: isRunning,
		},
	})
}

// StopCurrentAgent is the endpoint stop the currently running rendering agent
func (s Server) StopCurrentAgent(c *gin.Context) {
	if err := s.cndtr.StopCurrentAgent(); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, InvalidResponse{
			Success: false,
			Message: "No service is currently running",
		})
		return
	}
	c.JSON(http.StatusOK, ValidResponse[struct{}]{
		Success: true,
	})
}

// ListAgents will return a list of currently available rendering agents
func (s Server) ListAgents(c *gin.Context) {
	services := s.cndtr.ListAgents()
	c.JSON(200, ValidResponse[[]string]{
		Success:  true,
		Response: services,
	})
}

// UpdateConfig will update the configuration of the currently running agent
func (s Server) UpdateConfig(c *gin.Context) {
	var cfg map[string]string
	if c.ShouldBindJSON(&cfg) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: "invalid config, must be  object of string -> string",
		})
		return
	}
	if err := s.cndtr.UpdateConfig(cfg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: fmt.Sprintf("unable to update config: %s", err.Error()),
		})
		return
	}
}

// SetAgent sets the current agent rendering to the display
func (s Server) SetAgent(c *gin.Context) {
	var request SetAgentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Warningf("invalid set agent request: %s\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: "endpoint requires an agent name:",
		})
		return
	}
	if err := s.cndtr.InitNewAgent(request.Agent); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: fmt.Sprintf("unable to start agent %s: %s", request.Agent, err),
		})
		return
	}
	c.JSON(http.StatusAccepted, ValidResponse[struct{}]{
		Success: true,
	})
}
