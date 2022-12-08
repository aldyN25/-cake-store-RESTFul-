package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gitlab.com/cake-store-RESTFul/handler"
	service_manager "gitlab.com/cake-store-RESTFul/service-manager"
)

func v1(router *httprouter.Router, serviceManager service_manager.ServiceManager, log zerolog.Logger) {

	commonHttp := handler.NewCommonHttp()
	cakeHanlder := handler.NewCake(serviceManager.CakeService(), commonHttp, log)

	router.POST("/api/v1/cake", cakeHanlder.Create)
	router.POST("/api/v1/cake/json", cakeHanlder.CreateWithJSon)
	router.GET("/api/v1/cake", cakeHanlder.GetList)
	router.GET("/api/v1/cake/:id", cakeHanlder.GetDetail)
	router.PATCH("/api/v1/cake/:id", cakeHanlder.Update)
	router.DELETE("/api/v1/cake/:id", cakeHanlder.Delete)

}

func Run(viper *viper.Viper, serviceManager service_manager.ServiceManager, log zerolog.Logger) {

	config := viper.Sub("api")
	port := config.GetInt("port")
	address := fmt.Sprintf(":%d", port)
	_ = config.Get("environtment")

	router := httprouter.New()

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, _ interface{}) {
		res := handler.BaseResponse{Error: "internal server error", Data: nil}
		handler.NewCommonHttp().JSON(w, http.StatusInternalServerError, res)
	}

	v1(router, serviceManager, log)
	log.Fatal().Err(http.ListenAndServe(address, router)).Msg("service stop")

}
