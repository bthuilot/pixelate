package http

import (
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
)

type Server struct {
	cndtr  matrix.Conductor
	router *gin.Engine
}

type Options struct {
	Templates   fs.FS
	StaticFiles fs.FS
}

// NewServer will create a new HTTP Server
func NewServer(cndtr matrix.Conductor, options Options) (s *Server) {
	r := gin.Default()
	s = &Server{
		cndtr:  cndtr,
		router: r,
	}
	html, err := template.ParseFS(options.Templates, "*.tmpl")
	if err != nil {
		logrus.Fatalf("unable to read embded filesystem: %s", err)
	}
	r.SetHTMLTemplate(html)
	s.registerEndpoints(options.StaticFiles)
	cndtr.RegisterAgentEndpoints(r)
	return
}

// Run will start the HTTP server
func (s Server) Run() error {
	logrus.Info("Starting HTTP Server")
	return s.router.Run("0.0.0.0:8080") // listen and serve on localhost:8080
}
