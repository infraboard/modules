package impl

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/maudit/apps/event"
	"github.com/rs/zerolog"

	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	ioc.Controller().Registry(&impl{})
}

// 业务具体实现
type impl struct {
	// 继承模版
	ioc.ObjectImpl

	// 模块子Logger
	log *zerolog.Logger

	//
	col *mongo.Collection
}

// 对象名称
func (i *impl) Name() string {
	return event.AppName
}

// 初始化
func (i *impl) Init() error {
	// 对象
	i.log = log.Sub(i.Name())

	i.log.Debug().Msgf("database: %s", ioc_mongo.Get().Database)
	// 需要一个集合Collection
	i.col = ioc_mongo.DB().Collection("events")
	return nil
}
