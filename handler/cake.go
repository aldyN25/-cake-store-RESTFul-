package handler

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"gitlab.com/cake-store-RESTFul/service"
	cakeApi "gitlab.com/cake-store-RESTFul/service/cake"
	commonApi "gitlab.com/cake-store-RESTFul/service/common"
)

type Cake struct {
	cakeService service.Cake
	HttpSerializer
	Log zerolog.Logger
}

func NewCake(cakeService service.Cake, serializer HttpSerializer, log zerolog.Logger) *Cake {
	return &Cake{
		cakeService,
		serializer,
		log,
	}
}

func (c *Cake) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	reqBody := cakeApi.CreateRequest{}
	s, _, err := r.FormFile("image")
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: errors.New("image cannot be empty").Error(), Data: nil})
		return
	}

	reqBody.Image = s

	err = reqBody.ParseForm(r.Form)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	err = c.cakeService.Create(r.Context(), reqBody)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: "success"})
}

func (c *Cake) CreateWithJSon(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	reqBody := cakeApi.CreateRequestJSON{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&reqBody)
	if err != nil {
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	if err = reqBody.Validate(); err != nil {
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	b, err := base64.StdEncoding.DecodeString(reqBody.Image)
	if err != nil {
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
		return
	}
	buf := bytes.NewBuffer(b)

	err = c.cakeService.Create(r.Context(), cakeApi.CreateRequest{
		Title:       reqBody.Title,
		Rating:      reqBody.Rating,
		Description: reqBody.Description,
		Image:       buf,
	})

	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: "success"})
}

func (c *Cake) GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	urlQuery := r.URL.Query()
	req := cakeApi.GetListRequest{}
	req.ParseQuery(urlQuery)

	paginateReq := commonApi.PaginationRequest{}
	paginateReq.ParseQuery(urlQuery)

	res, pagination, err := c.cakeService.GetList(r.Context(), req, paginateReq)

	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: res, MetaData: pagination})
}

func (c *Cake) GetDetail(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	cakeIDStr := param.ByName("id")
	cakeID, err := strconv.Atoi(cakeIDStr)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
	}

	res, err := c.cakeService.GetDetail(r.Context(), cakeID)
	if err != nil && err == sql.ErrNoRows {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusNotFound, BaseResponse{Error: errors.New("data not found").Error(), Data: nil})
		return
	}

	if err != nil && err != sql.ErrNoRows {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: res})
}

func (c *Cake) Update(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	reqBody := cakeApi.UpdateRequest{}

	cakeIDStr := param.ByName("id")
	cakeID, err := strconv.Atoi(cakeIDStr)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
	}

	reqBody.ID = cakeID
	i, _, _ := r.FormFile("image")
	reqBody.Image = i

	err = reqBody.ParseForm(r.Form)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	err = c.cakeService.Update(r.Context(), reqBody)
	if err != nil {
		c.Log.Error().Err(err)
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: "success"})
}

func (c *Cake) Delete(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	cakeIDStr := param.ByName("id")
	cakeID, err := strconv.Atoi(cakeIDStr)
	if err != nil {
		c.Log.Error().Err(err)
		c.JSON(w, http.StatusBadRequest, BaseResponse{Error: err.Error(), Data: nil})
	}

	err = c.cakeService.Delete(r.Context(), cakeID)
	if err != nil {
		c.Log.Error().Err(err)
		c.JSON(w, http.StatusInternalServerError, BaseResponse{Error: err.Error(), Data: nil})
		return
	}

	c.JSON(w, http.StatusOK, BaseResponse{Error: nil, Data: "success"})
}
