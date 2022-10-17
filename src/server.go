package main

import (
	"embed"
	"fmt"
	"github.com/bthuilot/pixelate/matrix"
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

//go:embed web/static/js/* web/static/css/*
var staticFiles embed.FS

//go:embed web/templates/*
var templateFiles embed.FS

// ValidResponse represents a 200 success response
type ValidResponse[T interface{}] struct {
	Success  bool `json:"success"`
	Response T    `json:"response,omitempty"`
}

// InvalidResponse represents a failure from the server
type InvalidResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Server struct {
	cndtr  matrix.Conductor
	router *gin.Engine
}

// NewServer will create a new HTTP Server
func NewServer(cndtr matrix.Conductor) (s *Server) {
	r := gin.Default()
	s = &Server{
		cndtr:  cndtr,
		router: r,
	}
	html, err := template.ParseFS(templateFiles, "web/templates/*.tmpl")
	if err != nil {
		logrus.Fatalf("unable to read embded filesystem: %s", err)
	}
	r.SetHTMLTemplate(html)
	s.registerEndpoints()
	cndtr.RegisterAgentEndpoints(r)
	return
}

// Run will start the HTTP server
func (s Server) Run() error {
	logrus.Info("Starting HTTP Server")
	return s.router.Run("0.0.0.0:8080") // listen and serve on localhost:8080
}

// registerEndpoints will register all HTTP endpoints for the server
func (s Server) registerEndpoints() {
	logrus.Info("registering HTTP endpoints")
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	staticDir, err := fs.Sub(staticFiles, path.Join("web", "static"))
	if err != nil {
		panic("unable to traverse into static dir")
	}
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

// dashboard represents the struct to render the dashboard with
type dashboard struct {
	CurrentAgentRunning bool
	CurrentAgent        string
	Config              map[string]string
	Attributes          []template.HTML
	Agents              []string
}

// RenderDashboard is the endpoint to load and render the dashboard template
func (s Server) RenderDashboard(c *gin.Context) {
	name, cfg, attrs, running := s.cndtr.GetCurrentAgent()
	var html []template.HTML
	for _, attr := range attrs {
		html = append(html, template.HTML(attr.GetHTML()))
	}
	c.HTML(http.StatusOK, "dashboard.tmpl", dashboard{
		CurrentAgentRunning: running,
		CurrentAgent:        name,
		Config:              cfg,
		Attributes:          html,
		Agents:              s.cndtr.ListAgents(),
	})
}

// currentAgentResponse is the struct to represent the response from the HTTP server
// for the currently running agent
type currentAgentResponse struct {
	IsRunning bool              `json:"is_running"`
	ID        string            `json:"id"`
	Config    map[string]string `json:"config"`
}

// GetCurrentAgent is the endpoint to return the currently running agent
func (s Server) GetCurrentAgent(c *gin.Context) {
	id, cfg, _, isRunning := s.cndtr.GetCurrentAgent()
	c.JSON(http.StatusOK, ValidResponse[currentAgentResponse]{
		Success: true,
		Response: currentAgentResponse{
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
	}

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

// SetAgentRequest is a schema for a POST request to set the current agent
type setAgentRequest struct {
	Agent string `json:"agent"`
}

// SetAgent sets the current agent rendering to the display
func (s Server) SetAgent(c *gin.Context) {
	var request setAgentRequest
	if c.ShouldBindJSON(&request) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: "endpoint requires an agent name",
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
