package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_commonHttp_JSON(t *testing.T) {
	w := httptest.NewRecorder()
	type args struct {
		w      http.ResponseWriter
		status int
		res    BaseResponse
	}
	tests := []struct {
		name string
		c    *commonHttp
		args args
	}{
		{
			name: "success",
			c:    &commonHttp{},
			args: args{
				w:      w,
				res:    BaseResponse{},
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commonHttp{}
			c.JSON(tt.args.w, tt.args.status, tt.args.res)
		})
	}
}
