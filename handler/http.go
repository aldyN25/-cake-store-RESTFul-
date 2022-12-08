package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HttpSerializer interface {
	JSON(w http.ResponseWriter, status int, res BaseResponse)
}

type commonHttp struct{}

func NewCommonHttp() HttpSerializer {
	return new(commonHttp)
}

func (c *commonHttp) JSON(w http.ResponseWriter, status int, res BaseResponse) {

	jsonByte, err := json.Marshal(res)
	if err != nil {
		log.Err(err)
		panic(err)
	}

	w.WriteHeader(status)
	_, err = w.Write(jsonByte)
	if err != nil {
		log.Err(err)
		panic(err)
	}

}

type BaseResponse struct {
	Error    interface{} `json:"error"`
	Data     interface{} `json:"data"`
	MetaData interface{} `json:"meta_data,omitempty"`
}
