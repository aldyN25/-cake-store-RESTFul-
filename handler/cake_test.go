package handler

import (
	"database/sql"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	cakeAPi "gitlab.com/cake-store-RESTFul/service/cake"
	commonRes "gitlab.com/cake-store-RESTFul/service/common"
	mockService "gitlab.com/cake-store-RESTFul/service/mocks"
)

func TestCake_GetList(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cake", nil)
	tests := []struct {
		name       string
		w          http.ResponseWriter
		r          *http.Request
		in2        httprouter.Params
		beforeFunc func(s *mockService.MockCake)
	}{
		{
			name: "success",
			w:    res,
			r:    req,
			in2:  make(httprouter.Params, 1),
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetList(gomock.Any(), gomock.Any(), gomock.Any()).Return(cakeAPi.CakesResponse{}, commonRes.PaginationResponse{}, nil)
			},
		},
		{
			name: "error 400 (error parsing url query)",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "limit", Value: "invalid"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetList(gomock.Any(), gomock.Any(), gomock.Any()).Return(cakeAPi.CakesResponse{}, commonRes.PaginationResponse{}, nil)
			},
		},
		{
			name: "error 500",
			w:    res,
			r:    req,
			in2:  make(httprouter.Params, 1),
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetList(gomock.Any(), gomock.Any(), gomock.Any()).Return(cakeAPi.CakesResponse{}, commonRes.PaginationResponse{}, errors.New("foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			service := mockService.NewMockCake(ct)
			c := &Cake{
				cakeService:    service,
				HttpSerializer: new(commonHttp),
			}
			tt.beforeFunc(service)
			c.GetList(tt.w, tt.r, tt.in2)

		})
	}
}

func createImg() image.Image {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{X: 10, Y: 10}, Max: image.Point{X: 200, Y: 100}})
	img.Set(10, 10, color.White)
	return img
}

func TestCake_Create(t *testing.T) {

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("image", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		img := createImg()
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}

		err = writer.WriteField("rating", "1")
		if err != nil {
			t.Error(err)
		}

	}()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cake", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	tests := []struct {
		name       string
		w          http.ResponseWriter
		r          *http.Request
		in2        httprouter.Params
		beforeFunc func(s *mockService.MockCake)
	}{
		{
			name: "success",
			w:    res,
			r:    req,
			in2:  make(httprouter.Params, 0),
			beforeFunc: func(s *mockService.MockCake) {
				req.Form = url.Values{
					"rating":      []string{"1"},
					"title":       []string{"test"},
					"description": []string{"test"},
				}
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "error 500",
			w:    res,
			r:    req,
			in2:  make(httprouter.Params, 0),
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			service := mockService.NewMockCake(ct)

			c := &Cake{
				cakeService:    service,
				HttpSerializer: new(commonHttp),
			}
			tt.beforeFunc(service)
			c.Create(tt.w, tt.r, tt.in2)

		})
	}
}

func TestCake_CreateFailedParseFloat(t *testing.T) {

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("image", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		img := createImg()
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cake", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	ct := gomock.NewController(t)
	service := mockService.NewMockCake(ct)
	service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	c := &Cake{
		cakeService:    service,
		HttpSerializer: new(commonHttp),
	}
	c.Create(res, req, nil)
}

func TestCake_CreateFailedParseImg(t *testing.T) {

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("invalid", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		img := createImg()
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cake", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	ct := gomock.NewController(t)
	service := mockService.NewMockCake(ct)
	service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	c := &Cake{
		cakeService:    service,
		HttpSerializer: new(commonHttp),
	}
	c.Create(res, req, nil)
}

func TestCake_GetDetail(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cake/:id", nil)
	tests := []struct {
		name       string
		w          http.ResponseWriter
		r          *http.Request
		in2        httprouter.Params
		beforeFunc func(s *mockService.MockCake)
	}{
		{
			name: "success",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetDetail(gomock.Any(), 1).Return(cakeAPi.CakeResponse{}, nil)
			},
		},
		{
			name: "error 400 (inavlid param)",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "invalid"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetDetail(gomock.Any(), 0).Return(cakeAPi.CakeResponse{}, errors.New("foo"))
			},
		},
		{
			name: "error 404",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "invalid"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetDetail(gomock.Any(), 0).Return(cakeAPi.CakeResponse{}, sql.ErrNoRows)
			},
		},
		{
			name: "error 500",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().GetDetail(gomock.Any(), 1).Return(cakeAPi.CakeResponse{}, errors.New("foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			service := mockService.NewMockCake(ct)
			tt.beforeFunc(service)
			c := &Cake{
				cakeService:    service,
				HttpSerializer: new(commonHttp),
			}
			c.GetDetail(tt.w, tt.r, tt.in2)

		})
	}
}

func TestCake_Update(t *testing.T) {

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("image", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		img := createImg()
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}

		err = writer.WriteField("rating", "1")
		if err != nil {
			t.Error(err)
		}

	}()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cake/:id", pr)
	resF := httptest.NewRecorder()
	reqF := httptest.NewRequest(http.MethodPost, "/api/v1/cake/:id", nil)

	tests := []struct {
		name       string
		w          http.ResponseWriter
		r          *http.Request
		in2        httprouter.Params
		beforeFunc func(s *mockService.MockCake)
	}{
		{
			name: "success",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "error 500",
			w:    resF,
			r:    reqF,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				reqF.Form = url.Values{
					"rating": []string{"invalid"},
				}
				s.EXPECT().Update(gomock.Any(), cakeAPi.UpdateRequest{ID: 1}).Return(errors.New("foo"))
			},
		},
		{
			name: "error 400 (invalid param)",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "invalid"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Update(gomock.Any(), cakeAPi.UpdateRequest{ID: 0}).Return(errors.New("foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			service := mockService.NewMockCake(ct)
			tt.beforeFunc(service)
			c := NewCake(service, NewCommonHttp(), zerolog.Logger{})
			c.Update(tt.w, tt.r, tt.in2)
		})
	}
}

func TestCake_Delete(t *testing.T) {

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cake/:id", nil)
	tests := []struct {
		name       string
		w          http.ResponseWriter
		r          *http.Request
		in2        httprouter.Params
		beforeFunc func(s *mockService.MockCake)
	}{
		{
			name: "success",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
		},
		{
			name: "error 400",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "invalid"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Delete(gomock.Any(), 0).Return(errors.New("foo"))
			},
		},
		{
			name: "error 500",
			w:    res,
			r:    req,
			in2:  httprouter.Params{httprouter.Param{Key: "id", Value: "1"}},
			beforeFunc: func(s *mockService.MockCake) {
				s.EXPECT().Delete(gomock.Any(), 1).Return(errors.New("foo"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			service := mockService.NewMockCake(ct)
			c := NewCake(service, NewCommonHttp(), zerolog.Logger{})
			tt.beforeFunc(service)
			c.Delete(tt.w, tt.r, tt.in2)
		})
	}
}
