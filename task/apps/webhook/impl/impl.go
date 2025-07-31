package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/webhook"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&WebHookServiceImpl{
		WebhookQueue: "webhooks_run_queue",
	})
}

var _ webhook.Service = (*WebHookServiceImpl)(nil)

type WebHookServiceImpl struct {
	ioc.ObjectImpl

	// 当前这个消费者 配置的topic
	WebhookQueue string `toml:"webhook_topic" json:"webhook_topic" yaml:"webhook_topic"  env:"WEBHOOK_TOPIC"`

	cancelFn context.CancelFunc
	log      *zerolog.Logger
}

func (i *WebHookServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&webhook.WebHook{})
		if err != nil {
			return err
		}
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	i.cancelFn = cancelFn
	return i.RunWebHookConsumer(ctx)
}

func (i *WebHookServiceImpl) Name() string {
	return webhook.APP_NAME
}
