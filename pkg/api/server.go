package api

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/pkg/matrix"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type Server struct {
	selectedService Service
	matrix          *matrix.Service
	services        map[string]Service
	router          *gin.Engine
}

func CreateServer(allServices []Service) *Server {
	r := gin.Default()
	matrixService, err := matrix.CreateService()
	if err != nil {
		log.Fatalln(err)
	}
	s := createServices(r, matrixService, allServices)
	return &Server{
		selectedService: nil,
		services:        s,
		router:          r,
		matrix:          matrixService,
	}
}

func (s *Server) Run() {
	logging.InfoLogger.Println("Starting server... ")
	s.createEndpoints()
	logging.InfoLogger.Println("Spawning update loop")
	go func() {
		for {
			sleep := s.PerformTick()
			time.Sleep(sleep)
		}
	}()
	logging.InfoLogger.Println("Starting router")
	s.router.Run() // listen and serve on 0.0.0.0:8080
}

func (s *Server) PerformTick() time.Duration {
	logging.InfoLogger.Println("Polling services")
	if s.selectedService != nil {
		logging.InfoLogger.Printf("Performing tick for service %s\n", s.selectedService.GetID())
		err := s.selectedService.Tick()
		if err != nil {
			logging.WarningLogger.Printf("Unable to perform tick, error: %s\n", err)
		}
		return s.selectedService.RefreshDelay()
	} else {
		logging.InfoLogger.Println("No service selected, clearing screen")
		s.matrix.ClearScreen()
		return time.Second * 30
	}
}

func createServices(r *gin.Engine, matrix *matrix.Service, allServices []Service) (services map[string]Service) {
	matrix.Init()
	services = map[string]Service{}
	for _, service := range allServices {
		logging.InfoLogger.Printf("Initializing service %s\n", service.GetID())
		err := service.Init(matrix.Chan, r)
		if err != nil {
			logging.ErrorLogger.Printf("Unable to initialize serivce, error: %s\n", err)
			panic(err)
		}
		services[service.GetID()] = service
	}
	return
}
