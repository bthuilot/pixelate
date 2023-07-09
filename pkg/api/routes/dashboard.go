package routes

import (
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RenderDashboard is the endpoint to load and render the dashboard template
func RenderDashboard(d display.Display) gin.HandlerFunc {
	return func(c *gin.Context) {
		screen, name, running := d.CurrentScreen()
		c.HTML(http.StatusOK, "index.tmpl", struct {
			CurrentAgentRunning bool
			CurrentAgent        string
			Config              map[string]string
			Attributes          []display.HTMLAttributes
			Screens             []string
		}{
			CurrentAgentRunning: running,
			CurrentAgent:        name,
			Config:              screen.GetConfig(),
			Attributes:          screen.GetHTMLPage(),
			Screens:             d.GetScreens(),
		})
	}
}
