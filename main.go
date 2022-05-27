package main

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/spotify"
	"SpotifyDash/pkg/ticker"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	server := CreateServer()
	server.Run()
}

type Server struct {
	selectedService api.Service

	services []api.Service
	router   *gin.Engine
}

func CreateServer() Server {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	s := createServices(r)
	return Server{
		selectedService: s[1],
		services:        s,
		router:          r,
	}
}

func (s Server) Run() {
	go func() {
		fmt.Println("Here")
		for {
			if s.selectedService != nil {
				fmt.Println("tick")
				err := s.selectedService.Tick()
				if err != nil {
					fmt.Println(err)
				}
			}
			time.Sleep(time.Second * 5)
		}
		fmt.Println("WHtt end")
	}()
	s.router.Run() // listen and serve on 0.0.0.0:8080
}

func createServices(r *gin.Engine) []api.Service {
	matrixService, _ := matrix.CreateService()
	matrixService.Init()
	services := []api.Service{
		&spotify.Service{},
		&ticker.Service{},
	}
	for _, service := range services {
		err := service.Init(matrixService.Chan, r)
		if err != nil {
			panic(err)
		}
	}
	return services
}
