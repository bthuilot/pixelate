// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"github.com/bthuilot/pixelate/matrix"
	"github.com/bthuilot/pixelate/rendering"
	"github.com/bthuilot/pixelate/util"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load viper
	if err := util.InitConfig(); err != nil {
		panic(err)
	}
	// Create services
	logrus.Info("Creating renderers")
	rndrs := []rendering.Agent{
		rendering.NewSpotifyAgent(),
	}

	logrus.Info("Launching matrix service")
	// Create the conductor
	logrus.Info("Creating conductor")
	cndtr, err := matrix.NewConductor(rndrs)
	if err != nil {
		logrus.Error(err)
	}

	// Create and start webserver
	logrus.Info("Starting HTTP Server")
	server := NewServer(cndtr)
	if err := server.Run(); err != nil {
		logrus.Error(err)
	}
	logrus.Info("Server exited, shutting down agents")
	if err := cndtr.StopCurrentAgent(); err != nil {
		logrus.Error(err)
	}
}
