package httpsvr

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/pkg/util"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type ValidResponse[T interface{}] struct {
	Success  bool `json:"success"`
	Response T    `json:"response"`
}

type InvalidResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *Server) createEndpoints() {
	logging.InfoLogger.Println("Creating endpoint")
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.router.Static("/static", path.Join(util.GetDir(), "web/assets"))
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	// List all services
	s.router.GET("/services", s.ListServices())

	// Service
	s.router.GET("/service", s.GetCurrentService())
	s.router.POST("/service", s.SetService())
	s.router.DELETE("/service", s.StopCurrentService())

	// Config
	s.router.POST("/service/config", s.UpdateConfig())

	// Setup
	s.router.GET("/setup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup.tmpl", s.cndtr)
	})
}

type serviceResponse struct {
	IsRunning bool              `json:"is_running"`
	ID        string            `json:"id"`
	Config    map[string]string `json:"config"`
}

func (s *Server) GetCurrentService() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, cfg, isRunning := s.cndtr.GetCurrentService()
		c.JSON(http.StatusOK, ValidResponse[serviceResponse]{
			Success: true,
			Response: serviceResponse{
				ID:        id,
				Config:    cfg,
				IsRunning: isRunning,
			},
		})
	}
}

func (s *Server) StopCurrentService() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.cndtr.StopCurrentService()
	}
}

func (s *Server) ListServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.InfoLogger.Println("Retrieving services")
		services := s.cndtr.ListServices()
		c.JSON(200, ValidResponse[[]string]{
			Success:  true,
			Response: services,
		})
	}
}

func (s *Server) UpdateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cfg map[string]string
		err := c.ShouldBindJSON(&cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
				Success: false,
				Message: "invalid config, must be JSON object of string -> string",
			})
		}
	}
}

func (s *Server) SetService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service string
		err := c.ShouldBindJSON(&service)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
				Success: false,
				Message: "invalid config, must be JSON object of string -> string",
			})
		}
		err = s.cndtr.InitNewService(service)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, InvalidResponse{
				Success: false,
				Message: "unable to start service",
			})
		}
		c.JSON(http.StatusAccepted, ValidResponse[struct{}]{
			Success: true,
		})
	}
}
