package api

import (
	"github.com/gin-gonic/gin"
	"image"
)

type ConfigStore map[string]interface{}

type Service interface {
	GetConfig() ConfigStore
	SetConfig(config ConfigStore) error
	Init(matrixChan chan image.Image, engine *gin.Engine) error
	Tick() error
}
