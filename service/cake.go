package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rs/zerolog"
	"gitlab.com/cake-store-RESTFul/infra"
	"gitlab.com/cake-store-RESTFul/repo"
	cakeApi "gitlab.com/cake-store-RESTFul/service/cake"
	commonApi "gitlab.com/cake-store-RESTFul/service/common"
)

type Cake interface {
	Create(ctx context.Context, req cakeApi.CreateRequest) error
	GetList(ctx context.Context, req cakeApi.GetListRequest, paginateReq commonApi.PaginationRequest) (cakeApi.CakesResponse, commonApi.PaginationResponse, error)
	GetDetail(ctx context.Context, id int) (cakeApi.CakeResponse, error)
	Update(ctx context.Context, req cakeApi.UpdateRequest) error
	Delete(ctx context.Context, id int) error
}

type cake struct {
	cakeRepo   repo.Cake
	Log        zerolog.Logger
	Cloudinary infra.Cloudinary
}

func NewCake(cakeRepo repo.Cake, log zerolog.Logger, cl infra.Cloudinary) Cake {
	return &cake{
		cakeRepo:   cakeRepo,
		Log:        log,
		Cloudinary: cl,
	}
}

func (c *cake) Create(ctx context.Context, req cakeApi.CreateRequest) (err error) {

	r, err := c.Cloudinary.Upload(ctx, req.Image, uploader.UploadParams{})
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	cake := repo.CakeBaseModel{
		Title:       req.Title,
		Rating:      req.Rating,
		Description: req.Description,
		Image:       r.URL,
		CreatedAt:   time.Now().UTC(),
	}

	err = c.cakeRepo.Create(ctx, cake)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return err
	}

	return nil
}

func (c *cake) GetList(ctx context.Context, req cakeApi.GetListRequest, paginateReq commonApi.PaginationRequest) (res cakeApi.CakesResponse, pagination commonApi.PaginationResponse, err error) {

	offset := paginateReq.Limit * (paginateReq.Page - 1)
	cakes, err := c.cakeRepo.GetList(ctx, paginateReq.Limit, offset, req.Search, req.Sort, req.SortBy)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	for _, cake := range cakes {

		v := cakeApi.CakeResponse{
			ID:          cake.ID,
			Title:       cake.Title,
			Description: cake.Description,
			Image:       cake.Image,
			Rating:      cake.Rating,
			CreatedAt:   cake.CreatedAt,
		}

		if cake.UpdatedAt.Valid {
			v.UpdatedAt = &cake.UpdatedAt.Time
		}

		res = append(res, v)
	}

	if len(res) == 0 {
		res = make(cakeApi.CakesResponse, 0)
	}

	count, err := c.cakeRepo.CountCake(ctx, req.Search)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	pagination = commonApi.PaginationResponse{
		Page:  paginateReq.Page,
		Limit: paginateReq.Limit,
		Total: count,
	}

	return res, pagination, nil
}

func (c *cake) GetDetail(ctx context.Context, id int) (res cakeApi.CakeResponse, err error) {

	cake, err := c.cakeRepo.GetDetail(ctx, id)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return res, err
	}

	res = cakeApi.CakeResponse{
		ID:          cake.ID,
		Title:       cake.Title,
		Description: cake.Description,
		Image:       cake.Image,
		Rating:      cake.Rating,
		CreatedAt:   cake.CreatedAt,
	}

	if cake.UpdatedAt.Valid {
		res.UpdatedAt = &cake.UpdatedAt.Time
	}

	return res, nil
}

func (c *cake) Update(ctx context.Context, req cakeApi.UpdateRequest) (err error) {

	now := time.Now().UTC()
	iUrl := ""

	if req.Image != nil {
		r, err := c.Cloudinary.Upload(ctx, req.Image, uploader.UploadParams{})
		if err != nil {
			c.Log.Error().Msg(err.Error())
			return err
		}
		iUrl = r.URL
	}

	cake := repo.CakeBaseModel{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		UpdatedAt:   sql.NullTime{Time: now, Valid: true},
		Image:       iUrl,
	}

	err = c.cakeRepo.Update(ctx, cake)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return nil
}

func (c *cake) Delete(ctx context.Context, id int) (err error) {
	return c.cakeRepo.Delete(ctx, id)
}
