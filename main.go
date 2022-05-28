package main

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/spotify"
	"SpotifyDash/pkg/ticker"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	server := CreateServer()
	server.Run()
}

type Server struct {
	selectedService api.Service
	matrix          *matrix.Service
	services        []api.Service
	router          *gin.Engine
}

func CreateServer() Server {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	matrixService, err := matrix.CreateService()
	if err != nil {
		log.Fatalln(err)
	}
	s := createServices(r, matrixService)
	return Server{
		selectedService: s[1],
		services:        s,
		router:          r,
		matrix:          matrixService,
	}
}

func (s Server) Run() {
	go func() {
		fmt.Println("Here")
		sleep := time.Second * 30
		for {
			if s.selectedService != nil {
				fmt.Println("tick")
				err := s.selectedService.Tick()
				if err != nil {
					fmt.Println(err)
				}
				sleep = s.selectedService.RefreshDelay()
			} else {
				s.matrix.ClearScreen()
			}
			time.Sleep(sleep)
		}
		fmt.Println("WHtt end")
	}()
	s.router.Run() // listen and serve on 0.0.0.0:8080
}

func createServices(r *gin.Engine, matrix *matrix.Service) []api.Service {
	matrix.Init()
	services := []api.Service{
		&spotify.Service{},
		&ticker.Service{},
	}
	for _, service := range services {
		err := service.Init(matrix.Chan, r)
		if err != nil {
			panic(err)
		}
	}
	return services
}
