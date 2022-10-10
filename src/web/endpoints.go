package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

//go:embed static/js/* static/css/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

type ValidResponse[T interface{}] struct {
	Success  bool `json:"success"`
	Response T    `json:"response,omitempty"`
}

type InvalidResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s Server) registerEndpoints() {
	logrus.Info("registering HTTP endpoints")
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	staticDir, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("unable to draverse into static dir")
	}
	s.router.StaticFS("/static", http.FS(staticDir))
	s.router.GET("/", s.RenderDashboard)

	// List all services
	s.router.GET("/agents", s.ListServices)

	// Service
	s.router.GET("/agents/current", s.GetCurrentService)
	s.router.DELETE("/agents/current", s.StopCurrentService)
	s.router.POST("/agents/current", s.SetAgent)
	// Config
	s.router.POST("/agents/current/config", s.UpdateConfig)
}

type serviceResponse struct {
	IsRunning bool              `json:"is_running"`
	ID        string            `json:"id"`
	Config    map[string]string `json:"config"`
}

type Dashboard struct {
	CurrentAgentRunning bool
	CurrentAgent        string
	Config              map[string]string
	Attributes          []template.HTML
	Agents              []string
}

func (s Server) RenderDashboard(c *gin.Context) {
	name, cfg, attrs, running := s.cndtr.GetCurrentRenderer()
	var html []template.HTML
	for _, attr := range attrs {
		html = append(html, template.HTML(attr.GetHTML()))
	}
	c.HTML(http.StatusOK, "dashboard.tmpl", Dashboard{
		CurrentAgentRunning: running,
		CurrentAgent:        name,
		Config:              cfg,
		Attributes:          html,
		Agents:              s.cndtr.ListRenderers(),
	})
}

func (s Server) GetCurrentService(c *gin.Context) {
	id, cfg, _, isRunning := s.cndtr.GetCurrentRenderer()
	c.JSON(http.StatusOK, ValidResponse[serviceResponse]{
		Success: true,
		Response: serviceResponse{
			ID:        id,
			Config:    cfg,
			IsRunning: isRunning,
		},
	})
}

func (s Server) StopCurrentService(c *gin.Context) {
	if err := s.cndtr.StopCurrentRenderer(); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, InvalidResponse{
			Success: false,
			Message: "No service is currently running",
		})
	}

}

func (s Server) ListServices(c *gin.Context) {
	services := s.cndtr.ListRenderers()
	c.JSON(200, ValidResponse[[]string]{
		Success:  true,
		Response: services,
	})
}

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

type SetAgentRequest struct {
	Agent string `json:"agent"`
}

func (s Server) SetAgent(c *gin.Context) {
	var request SetAgentRequest
	if c.ShouldBindJSON(&request) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
			Success: false,
			Message: "endpoint requires an agent name",
		})
		return
	}
	if err := s.cndtr.InitNewRenderer(request.Agent); err != nil {
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
