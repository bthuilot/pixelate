package conductor

import (
	"SpotifyDash/pkg/matrix"
	"SpotifyDash/pkg/services"
	"fmt"
	"log"
)

type Conductor struct {
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
	return Conductor{
		services:       svcMap,
		setup:          nil,
		currentService: nil,
		matrix:         mtrx,
	}
}

func (c Conductor) ListServices() (result []services.ID) {
	for n := range c.services {
		result = append(result, n)
	}
	return
}

func (c Conductor) InitNewService(id services.ID) (err error) {
	svc, exist := c.services[id]
	if !exist {
		err = fmt.Errorf("cannot initialize no existant service %s", id)
	}
	_ = c.StopCurrentService // ignore the error, just want to stop a service if one is running
	c.setup = svc.Init(c.matrix.Chan)
	return nil
}

func (c Conductor) GetCurrentService() (id string, config services.Config, isRunning bool) {
	if c.currentService != nil {
		isRunning = true
		id = c.currentService.id
		config = c.currentService.config
	}
	return
}

func (c Conductor) UpdateConfig(newCfg services.Config) (err error) {
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

func (c Conductor) StopCurrentService() (err error) {
	if c.currentService != nil {
		c.currentService.channel <- services.Command{
			Code: services.Stop,
		}
	} else {
		err = fmt.Errorf("no service is currently running")
	}
	return
}
