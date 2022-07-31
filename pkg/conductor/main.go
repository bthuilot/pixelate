package conductor

import (
	"fmt"
	"log"

	"github.com/bthuilot/pixelate/pkg/matrix"
	"github.com/bthuilot/pixelate/pkg/services"
)

type Conductor interface {
	ListServices() []services.ID
	InitNewService(services.ID) error
	GetCurrentService() (string, services.Config, bool)
	UpdateConfig(services.Config) error
	GetSetup() (services.SetupPage, bool)
	StopCurrentService() error
}

type conductor struct {
	services       map[string]services.Service
	setup          services.SetupPage
	matrix         *matrix.Service
	currentService *runningService
}

type runningService struct {
	channel chan services.Command
	id      string
	config  services.Config
}

func SpawnConductor(mtrx *matrix.Service, svcs []services.Service) Conductor {
	svcMap := map[string]services.Service{}
	for _, s := range svcs {
		name := s.GetName()
		if _, e := svcMap[name]; e {
			log.Fatalf("service names should be unique, recieved %s", name)
		}
		svcMap[name] = s
	}
	return conductor{
		services:       svcMap,
		setup:          nil,
		currentService: nil,
		matrix:         mtrx,
	}
}

func (c conductor) GetSetup() (setup services.SetupPage, running bool) {
	setup = c.setup
	running = c.currentService != nil
	return
}

func (c conductor) ListServices() (result []services.ID) {
	for n := range c.services {
		result = append(result, n)
	}
	return
}

func (c conductor) InitNewService(id services.ID) (err error) {
	svc, exist := c.services[id]
	if !exist {
		err = fmt.Errorf("cannot initialize no existant service %s", id)
	}
	_ = c.StopCurrentService() // ignore the error, just want to stop a service if one is running
	c.setup = svc.Init(c.matrix.Chan)
	return nil
}

func (c conductor) GetCurrentService() (id string, config services.Config, isRunning bool) {
	if c.currentService != nil {
		isRunning = true
		id = c.currentService.id
		config = c.currentService.config
	}
	return
}

func (c conductor) UpdateConfig(newCfg services.Config) (err error) {
	if c.currentService != nil {
		c.currentService.channel <- services.Command{
			Code:   services.Update,
			Config: newCfg,
		}
	} else {
		err = fmt.Errorf("no service running")
	}
	return
}

func (c conductor) StopCurrentService() (err error) {
	if c.currentService != nil {
		c.currentService.channel <- services.Command{
			Code: services.Stop,
		}
	} else {
		err = fmt.Errorf("no service is currently running")
	}
	return
}
