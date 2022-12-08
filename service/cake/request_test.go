package cake

import (
	"mime/multipart"
	"net/url"
	"testing"
	"time"
)

func TestGetListRequest_ParseQuery(t *testing.T) {
	type fields struct {
		Search string
		Sort   string
		SortBy string
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
			name:    "success",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GetListRequest{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				SortBy: tt.fields.SortBy,
			}
			c.ParseQuery(tt.args.v)
		})
	}
}

func TestCreateRequest_ParseForm(t *testing.T) {
	type fields struct {
		Title       string
		Rating      float32
		Description string
		Image       multipart.File
		CreatedAt   time.Time
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
			name:   "success",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"1"},
				},
			},
			wantErr: false,
		},
		{
			name:   "invalid ratinng",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"invalid"},
				},
			},
			wantErr: true,
		},
		{
			name:   "error: empty title",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{""},
					"description": []string{"test"},
					"rating":      []string{"1"},
				},
			},
			wantErr: true,
		},
		{
			name:   "error: empty description",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{""},
					"rating":      []string{"1"},
				},
			},
			wantErr: true,
		},
		{
			name:   "error: empty rating",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{""},
				},
			},
			wantErr: true,
		},
		{
			name:   "error: negative rating",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"-1"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateRequest{
				Title:       tt.fields.Title,
				Rating:      tt.fields.Rating,
				Description: tt.fields.Description,
				Image:       tt.fields.Image,
				CreatedAt:   tt.fields.CreatedAt,
			}
			if err := c.ParseForm(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("CreateRequest.ParseForm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateRequest_ParseForm(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Rating      float32
		Description string
		Image       multipart.File
		UpdatedAt   time.Time
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
			name:   "success",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"1"},
				},
			},
			wantErr: false,
		},
		{
			name:   "invalid rating",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"invalid"},
				},
			},
			wantErr: true,
		},
		{
			name:   "error: negative rating",
			fields: fields{},
			args: args{
				url.Values{
					"title":       []string{"test"},
					"description": []string{"test"},
					"rating":      []string{"-1"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &UpdateRequest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Rating:      tt.fields.Rating,
				Description: tt.fields.Description,
				Image:       tt.fields.Image,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if err := c.ParseForm(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRequest.ParseForm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
