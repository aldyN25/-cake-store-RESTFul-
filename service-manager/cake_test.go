package service_manager

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"gitlab.com/cake-store-RESTFul/infra"
	"gitlab.com/cake-store-RESTFul/repo"
	"gitlab.com/cake-store-RESTFul/service"
	mockSvc "gitlab.com/cake-store-RESTFul/service/mocks"
)

func Test_serviceManager_CakeRepo(t *testing.T) {
	db, _, _ := sqlmock.New()
	cakeRepo := repo.NewCake(db, *zerolog.Ctx(context.Background()))

	infra := infra.Infra{
		Mysql: &infra.Mysql{
			MySQL: db,
		},
	}

	sm := NewServiceManager(&infra)

	if got := sm.CakeRepo(); !reflect.DeepEqual(got, got) {
		t.Errorf("serviceManager.CakeService() = %v, want %v", got, cakeRepo)
	}
}

func Test_serviceManager_CakeService(t *testing.T) {

	db, _, _ := sqlmock.New()
	log := *zerolog.Ctx(context.Background())
	cakeRepo := repo.NewCake(db, log)
	ctrl := gomock.NewController(t)
	mc := mockSvc.NewMockCloudinary(ctrl)
	cakeService := service.NewCake(cakeRepo, log, mc)

	infra := infra.Infra{
		Mysql: &infra.Mysql{
			MySQL: db,
		},
	}

	sm := NewServiceManager(&infra)

	if got := sm.CakeService(); !reflect.DeepEqual(got, got) {
		t.Errorf("serviceManager.CakeService() = %v, want %v", got, cakeService)
	}

}
