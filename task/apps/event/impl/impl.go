package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&EventServiceImpl{
		EventTopic: "events",
		GroupId:    "event_queue_workers",
	})
}

var _ event.Service = (*EventServiceImpl)(nil)

type EventServiceImpl struct {
	ioc.ObjectImpl
	log      *zerolog.Logger
	cancelFn context.CancelFunc

	// 当前这个消费者 配置的topic
	EventTopic string `toml:"event_topic" json:"event_topic" yaml:"event_topic"  env:"EVENT_TOPIC"`
	// 消费组Id
	GroupId string `toml:"group_id" json:"group_id" yaml:"group_id"  env:"GROUP_ID"`
}

func (i *EventServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&event.Event{})
		if err != nil {
			return err
		}
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	i.cancelFn = cancelFn

	// 启动消费者
	return i.SaveEvent(ctx)
}

func (i *EventServiceImpl) Name() string {
	return event.APP_NAME
}

func (i *EventServiceImpl) Close(ctx context.Context) {
	if i.cancelFn != nil {
		i.cancelFn()
	}
}
