package service_manager

import (
	"gitlab.com/cake-store-RESTFul/infra"
	"gitlab.com/cake-store-RESTFul/repo"
	"gitlab.com/cake-store-RESTFul/service"
)

type ServiceManager interface {
	// cake
	CakeRepo() repo.Cake
	CakeService() service.Cake
}

type serviceManager struct {
	infra *infra.Infra
}

func NewServiceManager(infra *infra.Infra) ServiceManager {
	return &serviceManager{
		infra: infra,
	}
}
