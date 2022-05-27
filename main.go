package main

import (
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/spotify"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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
		fmt.Println("Here")
		for {
			for _, s := range services {
				fmt.Println("tick")
				err := s.Tick()
				if err != nil {
					fmt.Println(err)
				}
			}
			time.Sleep(time.Second * 5)
		}
		fmt.Println("WHtt end")
	}()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func CreateServices(r *gin.Engine) []api.Service {
	matrixService, _ := matrix.CreateService()
	matrixService.Init()
	services := []api.Service{
		&spotify.Service{},
	}
	for _, service := range services {
		err := service.Init(matrixService.Chan, r)
		if err != nil {
			panic(err)
		}
	}
	return services
}
