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
	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load Logger
	logging.Init()

	server := api.CreateServer([]api.Service{
		&spotify.Service{},
		&ticker.Service{},
	})
	server.Run()
}
