package consumer

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/modules/maudit/apps/event"
)

// 读取消息，处理消息
func (c *consumer) Run(ctx context.Context) error {
	bus.GetService().Queue(ctx, c.AuditTopic, func(e *bus.Event) {
		ae := event.NewEvent()
		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err := ae.Load(e.Data)
		if err == nil {
			if err := event.GetService().SaveEvent(ctx, event.NewEventSet().Add(ae)); err != nil {
				c.log.Error().Msgf("save event error, %s", err)
			}
		}
	})
	return nil
}
