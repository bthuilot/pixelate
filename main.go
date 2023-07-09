// main is the starting point of the server
// will create the services, conductor and web server
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/bthuilot/pixelate/pkg/api"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/bthuilot/pixelate/pkg/display/screens"
	"github.com/bthuilot/pixelate/third_party/rgbmatrix"

	"github.com/bthuilot/pixelate/pkg/config"
	"github.com/bthuilot/pixelate/pkg/rendering"
	"github.com/sirupsen/logrus"
)

// newDisplayCanvas creates a new canvas that renders to the matrix display
func newDisplayCanvas() (*rgbmatrix.Canvas, error) {
	config := &rgbmatrix.DefaultConfig
	config.Cols = 64
	config.Rows = 64
	config.HardwareMapping = "adafruit-hat"
	config.Brightness = 50
	// create a new Matrix instance with the DefaultConfig
	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	return rgbmatrix.NewCanvas(matrix), err
}

func main() {
	var (
		cfg    config.ConfigFile
		err    error
		canvas *rgbmatrix.Canvas
	)
	// Load viper
	if cfg, err = config.InitConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err = config.InitLogger(cfg.Logging.Level, cfg.Logging.LogFile, cfg.Logging.UseSTDOUT); err != nil {
		logrus.Fatal(err)
	}

	// Load Embedded Files
	if err = initEmbed(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("creating screens")
	s := []display.Screen{
		// Load Spotify
		screens.NewSpotifyScreen(cfg),
		screens.NewWifiQRCode(cfg),
	}

	// Create services
	if err = rendering.LoadFonts(fonts); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Launching matrix service")
	if canvas, err = newDisplayCanvas(); err != nil {
		logrus.Fatal(err)
	}

	r := api.NewRouter(templateFiles, staticFiles)

	// Add custom routes
	logrus.Info("constructing display")
	d := display.NewDisplay(r.Group("/"), canvas, s)
	if err != nil {
		logrus.Fatal(err)
	}

	api.RegisterRoutes(r, d)

	// start webserver
	logrus.Info("Starting HTTP Server")
	if err = r.Run(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("server exited")
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
