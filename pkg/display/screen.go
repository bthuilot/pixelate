package display

import (
	"github.com/gin-gonic/gin"
	"image"
	"time"
)

type Screen interface {
	Init(r *gin.RouterGroup) (string, error)
	Render() (image.Image, time.Duration, error)
	GetConfig() map[string]string
	UpdateConfig(map[string]string) error
	GetHTMLPage() []HTMLAttributes
}
