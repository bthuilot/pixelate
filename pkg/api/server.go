package api

import (
	"SpotifyDash/internal/logging"
	"SpotifyDash/pkg/conductor"
	"SpotifyDash/pkg/util"
	"path"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cndtr  conductor.Conductor
	router *gin.Engine
}

func CreateServer(cndtr conductor.Conductor) *Server {
	r := gin.Default()
	r.LoadHTMLGlob(path.Join(util.GetDir(), "web", "templates", "*", "*.tmpl"))
	return &Server{
		cndtr:  cndtr,
		router: r,
	}
}

func (s *Server) Run() {
	logging.InfoLogger.Println("Starting server... ")
	s.createEndpoints()
	logging.InfoLogger.Println("Spawning update loop")
	logging.InfoLogger.Println("Starting router")
	s.router.Run("matrix.thuilot.io:80") // listen and serve on localhost:80
}
