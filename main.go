package main

import (
	"log"

	"github.com/bthuilot/pixelate/internal/logging"
	"github.com/bthuilot/pixelate/pkg/conductor"
	"github.com/bthuilot/pixelate/pkg/httpsvr"
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/bthuilot/pixelate/pkg/services"
)

func main() {
	// Load Logger
	if logging.Init() != nil {
		log.Fatal("Unable to open loggers")
	}

	svcs := []services.Service{
		services.Spotify{},
		services.Ticker{},
	}

	matrixService, err := matrix.CreateService()
	if err != nil {
		log.Fatalln(err)
	}
	cndtr := conductor.SpawnConductor(matrixService, svcs)
	server := httpsvr.CreateServer(cndtr)

	server.Run()
}
