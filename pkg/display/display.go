package display

import (
	"errors"
	"github.com/bthuilot/pixelate/third_party/rgbmatrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"image"
	"image/draw"
	"time"
)

type Display interface {
	SetScreen(name string) error
	GetScreens() []string
	CurrentScreen() (Screen, string, bool)
	ClearScreen() error
	SetScreenConfig(name string, cfg map[string]string) error
}

type display struct {
	t             *time.Timer
	canvas        *rgbmatrix.Canvas
	screens       map[string]Screen
	currentScreen string
}

func (d display) GetScreens() (names []string) {
	for n, _ := range d.screens {
		names = append(names, n)
	}
	return
}

func (d display) updateScreenFunc(name string, screen Screen) func() {
	return func() {
		img, dur, err := screen.Render()
		if err != nil {
			logrus.Errorf("error while rendering screen '%s': %s", name, err)
			// TODO(possible log to screen?)
		}
		d.draw(img)
		d.t.Reset(dur)
	}
}

var (
	ErrScreenNotFound  = errors.New("no screen found with that name")
	ErrNoInitialRender = errors.New("screen failed to perform initial render")
)

func (d display) CurrentScreen() (Screen, string, bool) {
	screen, exists := d.screens[d.currentScreen]
	return screen, d.currentScreen, exists && d.t != nil
}

func (d display) ClearScreen() error {
	if d.t != nil {
		d.t.Stop()
		d.t = nil
	}
	return d.canvas.Clear()
}

func (d display) SetScreen(name string) error {
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
	d.t = time.AfterFunc(initialDuration, d.updateScreenFunc(name, screen))
	return nil
}

func (d display) SetScreenConfig(name string, cfg map[string]string) error {
	screen, exists := d.screens[name]
	if !exists {
		logrus.Debugf("no screen found with name '%s', returning error", name)
		return ErrScreenNotFound
	}
	return screen.UpdateConfig(cfg)
	// TODO(make sure pointers work for updating config)
}

func (d display) draw(img image.Image) {
	draw.Draw(d.canvas, d.canvas.Bounds(), img, image.Point{}, draw.Src)
}

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
	return display{
		t:       nil,
		screens: initializedScreens,
		canvas:  canvas,
	}
}
