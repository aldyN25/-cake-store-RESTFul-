package common

import (
	"net/url"
	"testing"
)

func TestPaginationRequest_ParseQuery(t *testing.T) {

	val := url.Values{}
	val["limit"] = []string{"0"}
	val["page"] = []string{"0"}

	valErr := url.Values{}
	valErr["limit"] = []string{"error"}
	valErr["page"] = []string{"error"}

	type fields struct {
		Limit int
		Page  int
	}

	type args struct {
		v url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Limit: 0,
				Page:  0,
			},
			args: args{
				v: val,
			},
			wantErr: false,
		},
		{
			name: "success with zero value",
			fields: fields{
				Limit: 0,
				Page:  0,
			},
			args: args{
				v: make(url.Values),
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Limit: 0,
				Page:  1,
			},
			args: args{
				v: valErr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &PaginationRequest{
				Limit: tt.fields.Limit,
				Page:  tt.fields.Page,
			}
			c.ParseQuery(tt.args.v)
		})
	}
}
