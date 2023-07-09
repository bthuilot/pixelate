package display

import (
	"errors"
	"image"
	"image/draw"
	"time"

	"github.com/bthuilot/pixelate/third_party/rgbmatrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// ErrScreenNotFound is an error to represent when a screen is selected
	// that doesnt exist
	ErrScreenNotFound = errors.New("no screen found with that name")
	// ErrNoInitialRender is an error to represent when rendering a screen fails
	ErrNoInitialRender = errors.New("screen failed to perform initial render")
)

// Display is the interface for setting the current [Screen] of the
// matrix display
type Display interface {
	// SetScreen will set the current screen of the display
	// `name` must be on the of the strings returned from [GetScreens]
	SetScreen(name string) error
	// GetScreens will list all possible screen values
	GetScreens() []string
	// CurrentScreen will return the current [Screen] that is being rendered
	// and its formatted name
	CurrentScreen() (Screen, string, bool)
	// ClearScreen will stop the current [Screen] being rendered and reset it
	// to black
	ClearScreen() error
	// SetScreenConfig will update the configuration of a given [Screen]
	SetScreenConfig(name string, cfg map[string]string) error
}

// NewDisplay will construct a new [Display]. it will run the [Screen.Init] for each [Screen],
// allowing them to register custom HTTP routes in the [gin.RouterGroup]
func NewDisplay(r *gin.RouterGroup, canvas *rgbmatrix.Canvas, screens []Screen) Display {
	initializedScreens := map[string]Screen{}
	for _, s := range screens {
		name, err := s.Init(r)
		if err != nil {
			logrus.Errorf("unable to initalize screen: %s, skipping", err)
			continue
		}
		initializedScreens[name] = s
	}
	return &display{
		t:       nil,
		screens: initializedScreens,
		canvas:  canvas,
	}
}

type display struct {
	t             *time.Timer
	canvas        *rgbmatrix.Canvas
	screens       map[string]Screen
	currentScreen string
}

func (d *display) GetScreens() (names []string) {
	for n, _ := range d.screens {
		names = append(names, n)
	}
	return
}

func (d *display) updateScreenFunc(name string, screen Screen) func() {
	return func() {
		logrus.Debugf("rendering screen for %s", name)
		img, dur, err := screen.Render()
		if err != nil {
			logrus.Errorf("error while rendering screen '%s': %s", name, err)
			// TODO(possible log to screen?)
		}
		d.draw(img)
		if dur > 0 {
			d.t.Reset(dur)
		}
	}
}

func (d *display) CurrentScreen() (Screen, string, bool) {
	screen, exists := d.screens[d.currentScreen]
	return screen, d.currentScreen, exists
}

func (d *display) ClearScreen() error {
	logrus.Info("clearing screen")
	if d.t != nil {
		d.t.Stop()
		d.t = nil
	}
	return d.canvas.Clear()
}

func (d *display) SetScreen(name string) error {
	screen, exists := d.screens[name]
	if !exists {
		logrus.Debugf("no screen found with name '%s', returning error", name)
		return ErrScreenNotFound
	}
	if d.t != nil {
		_ = d.t.Stop()
		d.t = nil
	}
	initialImage, initialDuration, initErr := screen.Render()
	if initErr != nil {
		logrus.Errorf("unable to perform initial render for screen '%s': %s", name, initErr)
		return ErrNoInitialRender
	}
	d.draw(initialImage)
	if initialDuration > 0 {
		d.t = time.AfterFunc(initialDuration, d.updateScreenFunc(name, screen))
	}
	d.currentScreen = name
	return nil
}

func (d *display) SetScreenConfig(name string, cfg map[string]string) error {
	screen, exists := d.screens[name]
	if !exists {
		logrus.Debugf("no screen found with name '%s', returning error", name)
		return ErrScreenNotFound
	}
	return screen.UpdateConfig(cfg)
	// TODO(make sure pointers work for updating config)
}

func (d *display) draw(img image.Image) {
	logrus.Debugf("drawing image")
	draw.Draw(d.canvas, d.canvas.Bounds(), img, image.Point{}, draw.Src)
	if err := d.canvas.Render(); err != nil {
		logrus.Errorf("unable to render image: %s", err)
	}
}
