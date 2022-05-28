package main

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/spotify"
	"SpotifyDash/pkg/ticker"
	"github.com/joho/godotenv"
	"log"
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

	server := api.CreateServer([]api.Service{
		&spotify.Service{},
		&ticker.Service{},
	})

	server.Run()
}
