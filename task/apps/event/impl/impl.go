package impl

import (
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

	log *zerolog.Logger

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
	return nil
}

func (i *EventServiceImpl) Name() string {
	return event.APP_NAME
}
