// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"embed"
	"fmt"
	"github.com/bthuilot/pixelate/pkg/http"
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/bthuilot/pixelate/pkg/rendering"
	"github.com/bthuilot/pixelate/pkg/util"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path"
)

func main() {
	// Load viper
	if err := util.InitConfig(); err != nil {
		logrus.Fatal(err)
	}

	// Load Embedded Files
	if err := initEmbed(); err != nil {
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
		StaticFiles: staticFiles,
	})
	if err = server.Run(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server exited, shutting down agents")
	if err = cndtr.StopCurrentAgent(); err != nil {
		logrus.Fatal(err)
	}
}

/* Embedded files */

//go:embed assets/web/static/*
var staticFilesEmbed embed.FS
var staticFiles fs.FS

//go:embed assets/web/templates/*
var templateFilesEmbed embed.FS
var templateFiles fs.FS

//go:embed assets/fonts/*
var fontsEmbed embed.FS
var fonts fs.FS

func initEmbed() (err error) {
	// Fonts
	fontDir := path.Join("assets", "fonts")
	if os.Getenv("USE_FS") != "" {
		fonts = os.DirFS(fontDir)
	} else if fonts, err = fs.Sub(fontsEmbed, fontDir); err != nil {
		return fmt.Errorf("unable to trasverse fonts directory: %s", err)
	}

	// Static web content
	staticDir := path.Join("assets", "web", "static")
	if os.Getenv("USE_FS") != "" {
		staticFiles = os.DirFS(staticDir)
	} else if staticFiles, err = fs.Sub(staticFilesEmbed, staticDir); err != nil {
		return fmt.Errorf("unable to trasverse static files directory: %s", err)
	}

	// Web Templates
	templateDir := path.Join("assets", "web", "templates")
	if os.Getenv("USE_FS") != "" {
		templateFiles = os.DirFS(templateDir)
	} else if templateFiles, err = fs.Sub(templateFilesEmbed, templateDir); err != nil {
		return fmt.Errorf("unable to trasverse template files directory: %s", err)
	}
	return nil
}
