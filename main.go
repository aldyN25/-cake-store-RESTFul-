package main

import (
	"github.com/spf13/viper"
	"gitlab.com/cake-store-RESTFul/api"
	"gitlab.com/cake-store-RESTFul/infra"
	service_manager "gitlab.com/cake-store-RESTFul/service-manager"
)

func main() {

	appConfig := viper.New()
	appConfig.SetConfigName("app")
	appConfig.SetConfigType("toml")
	appConfig.AddConfigPath("./config/")

	err := appConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}

	infra := infra.NewInfra(appConfig)
	serviceManager := service_manager.NewServiceManager(infra)

	api.Run(appConfig, serviceManager, infra.Log)

}
