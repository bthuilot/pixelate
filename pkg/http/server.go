package http

import (
	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
)

// Server is an HTTP server for providing an API for controlling the matrix
type Server struct {
	// cndtr is the matrix.Conductor for interfacing with agents and the display
	cndtr matrix.Conductor
	// router is the HTTP router
	router *gin.Engine
}

// Options are options to provide
type Options struct {
	// Templates is a file system to read Go templates from
	Templates fs.FS
	// StaticFiles is the file system to read static files from
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
