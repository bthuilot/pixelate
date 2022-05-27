package api

import (
	"github.com/gin-gonic/gin"
	"image"
)

const (
	StringType = iota
	IntegerType
)

type ConfigValue struct {
	configType  int8
	configValue interface{}
}

type ConfigStore map[string]interface{}

type Service interface {
	// GetConfig is the retrieve the config for the
	GetConfig() ConfigStore
	SetConfig(config ConfigStore) error
	Init(matrixChan chan image.Image, engine *gin.Engine) error
	Tick() error
}
