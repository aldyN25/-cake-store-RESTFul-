package infra

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Infra struct {
	*Mysql
	Log        zerolog.Logger
	Cloudinary Cloudinary
}

func NewInfra(viper *viper.Viper) *Infra {
	return &Infra{
		Mysql:      newMysql(viper),
		Log:        newLogger(),
		Cloudinary: newCloudinary(viper),
	}
}
