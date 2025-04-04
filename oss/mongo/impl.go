package mongo

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/oss"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	log *zerolog.Logger
	db  *mongo.Database
	ioc.ObjectImpl
}

func (s *service) Init() error {
	s.db = ioc_mongo.DB()
	s.log = log.Sub("storage")
	return nil
}

func (s *service) Name() string {
	return oss.AppName
}
