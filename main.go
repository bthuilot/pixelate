package main

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/spotify"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	services := CreateServices(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go func() {
		for _, s := range services {
			s.Tick()
		}
	}()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func CreateServices(r *gin.Engine) []api.Service {
	matrixService, _ := matrix.CreateService()
	services := []api.Service{
		spotify.Service{},
	}
	for _, service := range services {
		err := service.Init(matrixService.Chan, r)
		if err != nil {
			panic(err)
		}
	}
	return services
}
