package main

import (
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/api"
	"SpotifyDash/pkg/types"
)

var services = [...]types.Service {
	Spo
}

func main() {
	api.Init()
	// spotify.Init()
	matrix.Init()
}
