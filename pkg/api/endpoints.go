package api

import (
	"SpotifyDash/internal/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) createEndpoints() {
	logging.InfoLogger.Println("Creating endpoint")
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Dashboard & Web pages
	s.router.Static("/static/", "/home/bryce/github/PiMatrix/web/assets")
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})
	s.router.GET("/config/:service", s.RenderConfigPage())
	// Services
	s.router.GET("/services", s.ListServices())
	s.router.GET("/service", s.GetCurrentService())
	s.router.POST("/service", s.SetService())
	s.router.DELETE("/service", s.RemoveService())

	// View Specific Services
	s.router.GET("/service/:service/config", s.GetServiceConfig())
	s.router.POST("/service/:service/config", s.UpdateServiceConfig())
}

func (s *Server) RenderConfigPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if serv, exists := s.services[c.Param("service")]; exists {
			c.HTML(http.StatusOK, "config.tmpl", gin.H{"serviceName": serv.GetID()})
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

func (s *Server) GetCurrentService() gin.HandlerFunc {
	return func(c *gin.Context) {
		if serv := s.selectedService; serv != nil {
			c.JSON(http.StatusOK, serv.GetID())
		} else {
			c.JSON(http.StatusOK, nil)
		}
	}
}

func (s *Server) RemoveService() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.selectedService = nil
		s.PerformTick()
		c.Status(http.StatusOK)
	}
}

func (s *Server) ListServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.InfoLogger.Println("Retrieving services")
		var services []string
		for _, service := range s.services {
			services = append(services, service.GetID())
		}
		response := gin.H{"services": services, "selected": nil}
		if s.selectedService != nil {
			response["selected"] = s.selectedService.GetID()
		}
		c.JSON(200, response)
	}
}

type UpdateServiceRequest struct {
	Service string `json:"service"`
}

func (s *Server) SetService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body UpdateServiceRequest
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		value := body.Service
		if service, exists := s.services[value]; exists {
			logging.InfoLogger.Printf("Setting service to %s\n", value)
			s.selectedService = service
			_ = s.PerformTick()
			c.Status(http.StatusOK)
			return
		}
		logging.WarningLogger.Printf("Bad value given: %s", value)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (s *Server) GetServiceConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := c.Param("service")
		if service, exist := s.services[serviceName]; exist {
			c.JSON(http.StatusOK, service.GetConfig())
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}

func (s *Server) UpdateServiceConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := c.Param("service")
		var config ConfigStore
		if err := c.ShouldBindJSON(&config); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if service, exist := s.services[serviceName]; exist {
			_ = service.SetConfig(config)
			s.PerformTick()
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}
