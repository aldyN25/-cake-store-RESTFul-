package infra

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/spf13/viper"
)

func (c *cld) Upload(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error) {
	return c.cloudinary.Upload.Upload(ctx, file, uploadParams)
}

type cld struct {
	cloudinary *cloudinary.Cloudinary
}
type Cloudinary interface {
	Upload(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)
}

func newCloudinary(viper *viper.Viper) Cloudinary {
	config := viper.Sub("cloudinary")
	cloudName := config.GetString("cloud_name")
	apiKey := config.GetString("api_key")
	secretKey := config.GetString("secret")

	cl, err := cloudinary.NewFromParams(cloudName, apiKey, secretKey)
	if err != nil {
		panic(err)
	}

	return &cld{
		cloudinary: cl,
	}
}
