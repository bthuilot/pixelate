// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"embed"
	"github.com/bthuilot/pixelate/pkg/http"
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/bthuilot/pixelate/pkg/rendering"
	"github.com/bthuilot/pixelate/pkg/util"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path"
)

//go:embed assets/web/static/*
var staticFilesEmbed embed.FS
var staticFiles, _ = fs.Sub(staticFilesEmbed, path.Join("assets", "web", "static"))

//go:embed assets/web/templates/*
var templateFilesEmbed embed.FS
var templateFiles, _ = fs.Sub(templateFilesEmbed, path.Join("assets", "web", "templates"))

//go:embed assets/fonts/*
var fontsEmbed embed.FS
var fonts, _ = fs.Sub(fontsEmbed, path.Join("assets", "fonts"))

func main() {
	// Load viper
	if err := util.InitConfig(); err != nil {
		logrus.Fatal(err)
	}
	// Create services
	logrus.Info("Creating renderers")
	rndrs := []rendering.Agent{
		rendering.NewSpotifyAgent(),
		rendering.NewTickerAgent(),
	}
	if err := rendering.LoadFonts(fonts); err != nil {
		logrus.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Launching matrix service")
	// Create the conductor
	logrus.Info("Creating conductor")
	cndtr, err := matrix.NewConductor(rndrs)
	if err != nil {
		logrus.Fatal(err)
	}

	// Create and start webserver
	logrus.Info("Starting HTTP Server")
	server := http.NewServer(cndtr, http.Options{
		Templates:   templateFiles,
		StaticFiles: os.DirFS("assets/web/static/"),
	})
	if err = server.Run(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server exited, shutting down agents")
	if err = cndtr.StopCurrentAgent(); err != nil {
		logrus.Fatal(err)
	}
}
