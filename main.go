// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"log"

	"github.com/bthuilot/pixelate/internal/logging"
	"github.com/bthuilot/pixelate/pkg/conductor"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/bthuilot/pixelate/pkg/httpsvr"
	"github.com/bthuilot/pixelate/pkg/matrix"
)

func main() {
	// Load Logger
	if logging.Init() != nil {
		log.Fatal("Unable to open loggers")
	}

	// Create services
	rndrs := []display.Renderer{
		display.Spotify{},
		display.Ticker{},
	}

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
