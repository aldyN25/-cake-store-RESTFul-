package service

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"reflect"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"gitlab.com/cake-store-RESTFul/repo"
	mockRepo "gitlab.com/cake-store-RESTFul/repo/mocks"
	cakeApi "gitlab.com/cake-store-RESTFul/service/cake"
	commonApi "gitlab.com/cake-store-RESTFul/service/common"
	mockSvc "gitlab.com/cake-store-RESTFul/service/mocks"
)

func Test_cake_Create(t *testing.T) {

	log := zerolog.Ctx(context.Background())
	ctx := context.Background()
	x := &multipart.FileHeader{}
	f, _ := x.Open()

	req := cakeApi.CreateRequest{
		Title:       "test",
		Rating:      1,
		Description: "test",
		Image:       f,
	}

	type args struct {
		ctx context.Context
		req cakeApi.CreateRequest
	}
	tests := []struct {
		name       string
		args       args
		beforeFunc func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary)
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, nil)
				m.EXPECT().Create(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error when call repository",
			args: args{
				ctx: ctx,
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, nil)
				m.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error when call upload file",
			args: args{
				ctx: ctx,
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := gomock.NewController(t)
			mc := mockSvc.NewMockCloudinary(m)
			cakeRepo := mockRepo.NewMockCake(m)
			tt.beforeFunc(cakeRepo, mc)
			c := NewCake(cakeRepo, *log, mc)

			if err := c.Create(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("cake.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cake_GetList(t *testing.T) {
	updateAt := sql.NullTime{Time: time.Now(), Valid: true}
	ctx := context.Background()
	req := cakeApi.GetListRequest{
		Search: "",
		Sort:   "id",
		SortBy: "asc",
	}

	res := cakeApi.CakesResponse{
		{
			ID:          1,
			Title:       "test",
			Description: "test",
			Image:       "",
			Rating:      1,
			UpdatedAt:   &updateAt.Time,
		},
	}

	output := []repo.CakeBaseModel{
		{
			ID:          1,
			Title:       "test",
			Description: "test",
			Image:       "",
			Rating:      1,
			UpdatedAt:   updateAt,
		},
	}

	page := commonApi.PaginationResponse{
		Page:  1,
		Limit: 10,
		Total: 10,
	}

	type args struct {
		ctx         context.Context
		req         cakeApi.GetListRequest
		paginateReq commonApi.PaginationRequest
	}
	tests := []struct {
		name           string
		args           args
		wantRes        cakeApi.CakesResponse
		wantPagination commonApi.PaginationResponse
		beforeFunc     func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary)
		wantErr        bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: req,
				paginateReq: commonApi.PaginationRequest{
					Limit: 10,
					Page:  1,
				},
			},
			wantRes:        res,
			wantPagination: page,
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().GetList(ctx, 10, 0, "", "id", "asc").Return(output, nil)
				m.EXPECT().CountCake(ctx, "").Return(10, nil)
			},
			wantErr: false,
		},
		{
			name: "error when call GetList function from repo",
			args: args{
				ctx: ctx,
				req: req,
				paginateReq: commonApi.PaginationRequest{
					Limit: 10,
					Page:  1,
				},
			},
			wantRes:        nil,
			wantPagination: commonApi.PaginationResponse{},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().GetList(ctx, 10, 0, "", "id", "asc").Return(nil, errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error when call Count()",
			args: args{
				ctx: ctx,
				req: req,
				paginateReq: commonApi.PaginationRequest{
					Limit: 10,
					Page:  1,
				},
			},
			wantRes:        res,
			wantPagination: commonApi.PaginationResponse{},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().GetList(ctx, 10, 0, "", "id", "asc").Return(output, nil)
				m.EXPECT().CountCake(ctx, "").Return(0, errors.New("foo"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			m := mockRepo.NewMockCake(ctrl)
			mc := mockSvc.NewMockCloudinary(ctrl)
			tt.beforeFunc(m, mc)
			c := NewCake(m, zerolog.Logger{}, mc)

			gotRes, gotPagination, err := c.GetList(tt.args.ctx, tt.args.req, tt.args.paginateReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("cake.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("cake.GetList() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("cake.GetList() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_cake_GetDetail(t *testing.T) {

	ctx := context.Background()
	updateAt := sql.NullTime{Time: time.Now(), Valid: true}
	output := repo.CakeBaseModel{
		ID:          1,
		Title:       "test",
		Description: "test",
		Rating:      1,
		UpdatedAt:   updateAt,
	}

	res := cakeApi.CakeResponse{
		ID:          1,
		Title:       "test",
		Description: "test",
		Rating:      1,
		UpdatedAt:   &updateAt.Time,
	}

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		wantRes    cakeApi.CakeResponse
		beforeFunc func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary)
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  1,
			},
			wantRes: res,
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().GetDetail(ctx, 1).Return(output, nil)
			},
			wantErr: false,
		},
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  1,
			},
			wantRes: cakeApi.CakeResponse{},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().GetDetail(ctx, 1).Return(repo.CakeBaseModel{}, errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := gomock.NewController(t)
			mc := mockSvc.NewMockCloudinary(m)
			cakeRepo := mockRepo.NewMockCake(m)
			tt.beforeFunc(cakeRepo, mc)
			c := NewCake(cakeRepo, zerolog.Logger{}, mc)

			gotRes, err := c.GetDetail(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("cake.GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("cake.GetDetail() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_cake_Update(t *testing.T) {
	ctx := context.Background()
	x := &multipart.FileHeader{}
	f, _ := x.Open()
	req := cakeApi.UpdateRequest{
		ID:          1,
		Title:       "test",
		Description: "test",
		Rating:      1,
		Image:       f,
	}

	type args struct {
		ctx context.Context
		req cakeApi.UpdateRequest
	}
	tests := []struct {
		name       string
		args       args
		beforeFunc func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary)
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				m.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error when call repo",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, nil)
				m.EXPECT().Update(context.Background(), gomock.Any()).Return(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error when upload file",
			args: args{
				ctx: context.Background(),
				req: req,
			},
			beforeFunc: func(m *mockRepo.MockCake, cl *mockSvc.MockCloudinary) {
				cl.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return(&uploader.UploadResult{}, errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			m := mockRepo.NewMockCake(ct)
			mc := mockSvc.NewMockCloudinary(ct)
			tt.beforeFunc(m, mc)
			c := NewCake(m, zerolog.Logger{}, mc)

			if err := c.Update(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("cake.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cake_Delete(t *testing.T) {

	m := gomock.NewController(t)
	mc := mockRepo.NewMockCake(m)
	mcd := mockSvc.NewMockCloudinary(m)
	mc.EXPECT().Delete(context.Background(), 1).Return(nil)
	c := NewCake(mc, zerolog.Logger{}, mcd)

	if err := c.Delete(context.Background(), 1); err != nil {
		t.Errorf("cake.Delete() error = %v", err)
	}
}
