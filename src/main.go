// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"github.com/bthuilot/pixelate/pkg/agents"
	"github.com/bthuilot/pixelate/pkg/conductor"
	"github.com/bthuilot/pixelate/pkg/httpsvr"
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	// Load Logger

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
	// Create the conductor
	cndtr := conductor.SpawnConductor(matrixService, rndrs)

	// Create and start webserver
	server := httpsvr.CreateServer(cndtr)
	server.Run()
}
