package services

import (
	"image"
	"time"
)

type ID = string

type Config map[string]string

type Service interface {
	GetName() ID
	Run(Config) chan Command
	GetDefaultConfig() Config
	GetTickInterval() time.Duration
	Init(chan image.Image) SetupPage
}

type ConfigType = int

type SetupPage []Attribute

type CommandCode = int

const (
	Stop CommandCode = iota
	Update
	Tick
)

type Command struct {
	Code   CommandCode
	Config Config
}
