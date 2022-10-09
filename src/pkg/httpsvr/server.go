package httpsvr

import (
	"github.com/bthuilot/pixelate/pkg/conductor"
	"github.com/bthuilot/pixelate/pkg/util"
	"github.com/sirupsen/logrus"
	"path"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cndtr  conductor.Conductor
	router *gin.Engine
}

func CreateServer(cndtr conductor.Conductor) (s *Server) {
	r := gin.Default()
	r.LoadHTMLGlob(path.Join(util.GetDir(), "web", "templates", "*.tmpl"))
	s = &Server{
		cndtr:  cndtr,
		router: r,
	}
	s.registerEndpoints()
	cndtr.RegisterAgentEndpoints(r)
	return
}

func (s Server) Run() {
	logrus.Info("Starting HTTP Server")
	s.router.Run("0.0.0.0:8080") // listen and serve on localhost:8080
}
