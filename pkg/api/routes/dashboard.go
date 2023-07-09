package routes

import (
	"net/http"

	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/gin-gonic/gin"
)

// RenderDashboard is the endpoint to load and render the dashboard template
func RenderDashboard(d display.Display) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			cfg map[string]string
			// This is done to allow HTML to be rendered 'unsafely'
			attrs []display.HTMLAttribute
		)
		screen, name, running := d.CurrentScreen()
		if running {
			cfg = screen.GetConfig()
			attrs = screen.GetHTMLPage()
		}

		c.HTML(http.StatusOK, "index.tmpl", struct {
			CurrentAgentRunning bool
			CurrentAgent        string
			Config              map[string]string
			Attributes          []display.HTMLAttribute
			Screens             []string
		}{
			CurrentAgentRunning: running,
			CurrentAgent:        name,
			Config:              cfg,
			Attributes:          attrs,
			Screens:             d.GetScreens(),
		})
	}
}
