package main

import (
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/server"
	"SpotifyDash/pkg/spotify"
)


func main() {
	server.Init()
	spotify.Init()
	matrix.Init()
}




