package core

import (
	"fmt"
	"github.com/kardianos/service"
	"log"
)

type Program struct{}

var program Program
var svc service.Service

func InitializeService(serviceConf *service.Config) {
	program = Program{}
	svcConfig := serviceConf

	svcLocal, err := service.New(&program, svcConfig)
	svc = svcLocal
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func GetServiceInstance(serviceConf *service.Config) service.Service {
	InitializeService(serviceConf)
	return svc
}

func (p *Program) run() {}

func (p *Program) Start(s service.Service) error {
	fmt.Println("Now Service Start")
	return nil
}

func (p *Program) Stop(s service.Service) error {
	return nil
}

func (p *Program) Restart(s service.Service) error {
	fmt.Println("Now Service Restart")
	return nil
}