package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/audit"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&AuditServiceImpl{})
}

// 业务具体实现
type AuditServiceImpl struct {
	// 继承模版
	ioc.ObjectImpl

	// 当前这个消费者 配置的topic
	AuditTopic string `toml:"topic" json:"topic" yaml:"topic"  env:"TOPIC"`

	cancelFn context.CancelFunc
	// 模块子Logger
	log *zerolog.Logger
}

// 对象名称
func (i *AuditServiceImpl) Name() string {
	return audit.AppName
}

func (i *AuditServiceImpl) Priority() int {
	return audit.PRIORITY
}

// 初始化
func (i *AuditServiceImpl) Init() error {
	// 对象
	i.log = log.Sub(i.Name())

	ctx, cancelFn := context.WithCancel(context.Background())
	i.cancelFn = cancelFn
	return i.RunSaveEventConsumer(ctx)
}
