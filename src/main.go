// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"log"

	"github.com/bthuilot/pixelate/conductor"
	"github.com/bthuilot/pixelate/matrix"
	"github.com/bthuilot/pixelate/rendering"
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
	rndrs := []rendering.Agent{
		rendering.NewSpotifyAgent(),
	}

	logrus.Info("Launching matrix service")
	// Create the matrix service
	mtrx, err := matrix.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		mtrx.ClearScreen()
	}()
	// Create the conductor
	cndtr := conductor.New(mtrx, rndrs)

	// Create and start webserver
	server := web.NewServer(cndtr)
	if err := server.Run(); err != nil {
		logrus.Error(err)
	}
	shutdown(mtrx, cndtr)
}

func shutdown(mtrx *matrix.Service, cndtr conductor.Conductor) {
	if err := cndtr.StopCurrentRenderer(); err != nil {
		logrus.Error(err)
	}
	mtrx.Exit <- struct{}{}
}
