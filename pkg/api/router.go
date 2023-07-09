package api

import (
	"github.com/bthuilot/pixelate/pkg/api/routes"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
	"net/http"
)

func NewRouter(templateFS fs.FS, staticFS fs.FS) *gin.Engine {
	// create webserver
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFS(templateFS, "*.tmpl")))
	r.StaticFS("/static", http.FS(staticFS))
	return r
}

func RegisterRoutes(r *gin.Engine, d display.Display) {

	logrus.Info("registering HTTP endpoints")

	r.GET("/health", routes.Health())

	/* HTML Pages */
	// Dashboard
	r.GET("/", routes.RenderDashboard(d))
	/* Agents */
	// All Agents
	//s.router.GET("/screens", s.ListAgents)
	// Current Agent
	//s.router.GET("/agents/current", s.GetCurrentAgent)
	r.DELETE("/screens/current", routes.ClearScreen(d))
	r.POST("/screens/current", routes.SetScreen(d))
	// Config
	r.POST("/screens/current/config", routes.UpdateScreenConfig(d))
}

//// GetCurrentAgent is the endpoint to return the currently running agent
//func (s Server) GetCurrentAgent(c *gin.Context) {
//	id, cfg, _, isRunning := s.cndtr.GetCurrentAgent()
//	c.JSON(http.StatusOK, ValidResponse[CurrentAgentResponse]{
//		Success: true,
//		Response: CurrentAgentResponse{
//			ID:        id,
//			Config:    cfg,
//			IsRunning: isRunning,
//		},
//	})
//}

//// ListAgents will return a list of currently available rendering agents
//func (s Server) ListAgents(c *gin.Context) {
//	services := s.cndtr.ListAgents()
//	c.JSON(200, responses.ValidResponse[[]string]{
//		Success:  true,
//		Response: services,
//	})
//}
