package service_manager

import (
	"sync"

	"gitlab.com/cake-store-RESTFul/repo"
	"gitlab.com/cake-store-RESTFul/service"
)

var (
	cakeRepo        repo.Cake
	cakeRepoOnce    sync.Once
	cakeService     service.Cake
	cakeServiceOnce sync.Once
)

func (s *serviceManager) CakeRepo() repo.Cake {
	cakeRepoOnce.Do(func() {
		cakeRepo = repo.NewCake(s.infra.MySQL, s.infra.Log)
	})
	return cakeRepo
}

func (s *serviceManager) CakeService() service.Cake {
	cakeServiceOnce.Do(func() {
		cakeService = service.NewCake(s.CakeRepo(), s.infra.Log, s.infra.Cloudinary)
	})
	return cakeService
}
