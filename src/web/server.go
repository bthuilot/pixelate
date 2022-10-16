package web

import (
	"html/template"

	"github.com/bthuilot/pixelate/conductor"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cndtr  conductor.Conductor
	router *gin.Engine
}

// NewServer will create a new HTTP Server
func NewServer(cndtr conductor.Conductor) (s *Server) {
	r := gin.Default()
	s = &Server{
		cndtr:  cndtr,
		router: r,
	}
	html, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		logrus.Fatalf("unable to read embded filesystem: %s", err)
	}
	r.SetHTMLTemplate(html)
	s.registerEndpoints()
	cndtr.RegisterAgentEndpoints(r)
	return
}

func (s Server) Run() error {
	logrus.Info("Starting HTTP Server")
	return s.router.Run("0.0.0.0:8080") // listen and serve on localhost:8080
}
