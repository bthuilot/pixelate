package main

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/pkg/conductor"
	"SpotifyDash/pkg/httpsvr"
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	// Load env
	if godotenv.Load("secrets.env") != nil {
		log.Fatal("Error loading .env file")
	}

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
