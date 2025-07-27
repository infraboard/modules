package consumer

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/maudit/apps/event"
	"github.com/rs/zerolog"

	ioc_mongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
)

func init() {
	ioc.Controller().Registry(&consumer{
		GroupId:    "maudit",
		AuditTopic: "maudit",
		ctx:        context.Background(),
	})
}

// 业务具体实现
type consumer struct {
	// 继承模版
	ioc.ObjectImpl

	// 模块子Logger
	log *zerolog.Logger

	// 允许时上下文
	ctx context.Context

	// 消费组Id
	GroupId string `toml:"group_id" json:"group_id" yaml:"group_id"  env:"GROUP_ID"`
	// 当前这个消费者 配置的topic
	AuditTopic string `toml:"topic" json:"topic" yaml:"topic"  env:"TOPIC"`
}

// 对象名称
func (i *consumer) Name() string {
	return "maudit_consumer"
}

func (i *consumer) Priority() int {
	return event.PRIORITY - 1
}

// 初始化
func (i *consumer) Init() error {
	// 对象
	i.log = log.Sub(i.Name())
	i.log.Debug().Msgf("database: %s", ioc_mongo.Get().Database)
	i.Run(i.ctx)
	return nil
}

func (i *consumer) Close(ctx context.Context) {
	i.ctx.Done()
}
