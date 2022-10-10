// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"log"

	"github.com/bthuilot/pixelate/agents"
	"github.com/bthuilot/pixelate/conductor"
	"github.com/bthuilot/pixelate/matrix"
	"github.com/bthuilot/pixelate/util"
	"github.com/bthuilot/pixelate/web"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load viper
	if err := util.InitConfig(); err != nil {
		panic(err)
	}
	// Create services
	logrus.Info("Creating renderers")
	rndrs := []agents.Renderer{
		agents.NewSpotify(),
	}

	logrus.Info("Launching matrix service")
	// Create the matrix service
	matrixService, err := matrix.CreateService()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		matrixService.ClearScreen()
	}()
	// Create the conductor
	cndtr := conductor.SpawnConductor(matrixService, rndrs)

	// Create and start webserver
	server := web.CreateServer(cndtr)
	server.Run()
}
